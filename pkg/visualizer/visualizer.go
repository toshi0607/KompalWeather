package visualizer

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/sclevine/agouti"
	"github.com/toshi0607/kompal-weather/internal/config"
	"github.com/toshi0607/kompal-weather/pkg/gcs"
	"github.com/toshi0607/kompal-weather/pkg/logger"
)

// Visualizer represents visualizer
type Visualizer struct {
	config *config.VisualizerConfig
	gcs    *gcs.GCS
	log    logger.Logger
}

const (
	maleFileName         = "男湯サウナ.png"
	femaleFileName       = "女湯サウナ.png"
	lastPagePNGFileName  = "last-page.png"
	lastPageHTMLFileName = "last-page.html"
)

// New builds new Visualizer
func New(c *config.VisualizerConfig, g *gcs.GCS, l logger.Logger) (*Visualizer, error) {
	return &Visualizer{config: c, gcs: g, log: l}, nil
}

// Save saves male and female report files in GCS
func (v Visualizer) Save(ctx context.Context, rt ReportType) ([]string, error) {
	hasMale, hasFemale, err := v.hasFile(ctx, rt)
	if err != nil {
		return nil, fmt.Errorf("failed to check file existence: %v", err)
	}
	if hasMale && hasFemale {
		v.log.Info("male & female files already exist in GCS")
		return nil, nil
	}

	localPath, err := ioutil.TempDir("", fmt.Sprintf("%v", time.Now().Unix()))
	v.log.Info("localPath: %s", localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create tmp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(localPath); err != nil {
			v.log.Error("failed to remove tmp files", err)
		}
	}()

	driver, err := v.initDriver(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init driver: %v", err)
	}
	defer func() {
		if err := driver.Stop(); err != nil {
			v.log.Error("failed to stop driver", err)
		}
	}()
	page, err := driver.NewPage()
	if err != nil {
		return nil, fmt.Errorf("failed to open new page: %v", err)
	}
	defer func() {
		v.log.Info("collecting logs...")
		if err := page.Screenshot(fmt.Sprintf("%s/%s", localPath, lastPagePNGFileName)); err != nil {
			v.log.Error("failed to take screenshot, error:", err)
		}
		html, err := page.HTML()
		if err != nil {
			v.log.Error("failed to get HTML string, error:", err)
		}
		hb := []byte(html)
		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", localPath, lastPageHTMLFileName), hb, 0644); err != nil {
			v.log.Error("failed to save HTML, error:", err)
		}

		if _, err := v.uploadFiles(ctx, localPath, lastPagePNGFileName, ""); err != nil {
			v.log.Error("failed to upload PNG, error:", err)
		}
		if _, err := v.uploadFiles(ctx, localPath, lastPageHTMLFileName, ""); err != nil {
			v.log.Error("failed to upload HTML, error:", err)
		}
	}()

	lp, err := newLoginPage(page, v.log)
	if err != nil {
		return nil, fmt.Errorf("failed to open login page: %v", err)
	}

	loggedIn, err := lp.login(v.config.Mail, v.config.PW)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	mp, err := newMonitoringPage(loggedIn, v.log)
	if err != nil {
		return nil, fmt.Errorf("failed to open monitoring page: %v", err)
	}

	// When the login attempt is considered to be suspicious
	// if err := mp.page.FindByID("Email").Fill(v.config.Mail); err != nil {
	// 	 return "", fmt.Errorf("failed to fill login input: %v", err)
	// }
	// if err := mp.page.FindByID("next").Click(); err != nil {
	// 	 return "", fmt.Errorf("failed to click ID next button: %v", err)
	// }
	// time.Sleep(15 * time.Second)
	//
	// if err := mp.page.FindByID("password").Fill(v.config.PW); err != nil {
	// 	 return "", fmt.Errorf("failed to fill pw input: %v", err)
	// }
	// if err := mp.page.FindByID("submit").Click(); err != nil {
	//   return "", fmt.Errorf("failed to click pw submit button: %v", err)
	// }
	// time.Sleep(15 * time.Second)
	//
	// if err := mp.page.FindByID("submit").Click(); err != nil {
	//   return "", fmt.Errorf("failed to click pw submit button: %v", err)
	// }
	// time.Sleep(30 * time.Second)
	// It's necessary to choose the same number on my smartphone as shown in the screen on cloud...
	// When this happens often, I need to implement slack notification with screenshot.

	if err := mp.download(rt); err != nil {
		return nil, fmt.Errorf("failed to download files: %v", err)
	}

	var files []string
	if !hasMale {
		file, err := v.uploadFiles(ctx, localPath, maleFileName, rt)
		if err != nil {
			return nil, fmt.Errorf("failed to update male file: %v", err)
		}
		files = append(files, file)
	}
	if !hasFemale {
		file, err := v.uploadFiles(ctx, localPath, femaleFileName, rt)
		if err != nil {
			return nil, fmt.Errorf("failed to update female file: %v", err)
		}
		files = append(files, file)
	}

	return files, nil
}

func (v Visualizer) hasFile(ctx context.Context, rt ReportType) (bool, bool, error) {
	malePath, err := v.objectPath(maleFileName, rt)
	if err != nil {
		return false, false, fmt.Errorf("failed to build male object path: %v", err)
	}
	hasMale, err := v.gcs.HasObject(ctx, malePath)
	if err != nil {
		return false, false, fmt.Errorf("failed to check male object: %v", err)
	}
	femalePath, err := v.objectPath(femaleFileName, rt)
	if err != nil {
		return false, false, fmt.Errorf("failed to build female object path: %v", err)
	}
	hasFemale, err := v.gcs.HasObject(ctx, femalePath)
	if err != nil {
		return false, false, fmt.Errorf("failed to check female object: %v", err)
	}
	return hasMale, hasFemale, nil
}

func (v Visualizer) initDriver(localPath string) (*agouti.WebDriver, error) {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("prefs", map[string]interface{}{
			"download.default_directory": localPath,
		}),
		agouti.ChromeOptions("args", []string{
			"--headless",             // headless mode
			"--window-size=1280,800", // Size of window
			"--no-sandbox",           // Sandbox requires namespace permissions that we don't have on a container
			"--disable-gpu",          // There is no GPU on our Ubuntu box
			"--lang=ja",              // Language
		}),
		agouti.Debug,
	)

	if err := driver.Start(); err != nil {
		return nil, fmt.Errorf("failed to start driver: %v", err)
	}
	return driver, nil
}

func (v Visualizer) uploadFiles(ctx context.Context, localPath, fileName string, rt ReportType) (string, error) {
	f, err := os.Open(localPath + "/" + fileName)
	if err != nil {
		return "", fmt.Errorf("failed to open local file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			v.log.Error("failed to close local file", err)
		}
	}()
	fileInfo, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}

	b := make([]byte, fileInfo.Size())
	buffer := bufio.NewReader(f)
	_, err = buffer.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to read local file to buffer: %v", err)
	}

	op, err := v.objectPath(fileName, rt)
	if err != nil {
		return "", fmt.Errorf("failed to build object path: %v", err)
	}
	if err := v.gcs.Put(ctx, b, op); err != nil {
		return "", err
	}
	v.log.Info("file uploaded! bucket: %s, objectPath: %s", v.config.BucketName, op)

	return op, nil
}

// objectPath returns full path for GCS
// Example:
//   daily:   daily/2020-09-09-2020-09-09-male.png
//   weekly:  weekly/2020-12-28-2021-01-03-female.png
//   monthly: monthly/2020-12-01-2020-12-31-male.png
//   logs:    logs/1599983507/last-page.png
func (v Visualizer) objectPath(fileName string, rt ReportType) (string, error) {
	if rt == "" {
		return fmt.Sprintf("logs/%v/%s", time.Now().Unix(), fileName), nil
	}
	var gender string
	if fileName == maleFileName {
		gender = "male"
	} else {
		gender = "female"
	}

	periodStr, err := rt.reportPeriod()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s-%s.png", rt, periodStr, gender), nil
}
