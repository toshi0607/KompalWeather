package visualizer

import (
	"fmt"
	"time"

	"github.com/sclevine/agouti"
)

const (
	loginPageURL      = "https://stackoverflow.com/users/signup"
	monitoringPageURL = "https://console.cloud.google.com/monitoring/dashboards/custom/3717473642519746493"
)

type loginPage struct {
	page *agouti.Page
}

func newLoginPage(p *agouti.Page) (*loginPage, error) {
	if err := p.Navigate(loginPageURL); err != nil {
		return nil, err
	}
	time.Sleep(5 * time.Second)
	return &loginPage{page: p}, nil
}

//
// Page components
//
func (p loginPage) googleOAuthButton() *agouti.Selection {
	return p.page.FindByXPath("//*[@id=\"openid-buttons\"]/button[1]")
}

func (p loginPage) loginInput() *agouti.Selection {
	return p.page.FindByID("Email")
}

func (p loginPage) idNextButton() *agouti.Selection {
	return p.page.FindByID("next")
}

func (p loginPage) passwordInput() *agouti.Selection {
	return p.page.FindByID("password")
}

func (p loginPage) passwordNextButton() *agouti.Selection {
	return p.page.FindByID("submit")
}

//
// Page actions
//
func (p loginPage) login(id, pw string) (*agouti.Page, error) {
	if err := p.googleOAuthButton().Click(); err != nil {
		return nil, fmt.Errorf("failed to click google OAuth button: %v", err)
	}
	time.Sleep(15 * time.Second)

	if err := p.loginInput().Fill(id); err != nil {
		return nil, fmt.Errorf("failed to fill login input: %v", err)
	}
	if err := p.idNextButton().Click(); err != nil {
		return nil, fmt.Errorf("failed to click ID next button: %v", err)
	}
	time.Sleep(10 * time.Second)
	if err := p.passwordInput().Fill(pw); err != nil {
		return nil, fmt.Errorf("failed to fill password input: %v", err)
	}
	if err := p.passwordNextButton().Click(); err != nil {
		return nil, fmt.Errorf("failed to click password next button: %v", err)
	}
	time.Sleep(30 * time.Second)

	return p.page, nil
}

type monitoringPage struct {
	page *agouti.Page
}

func newMonitoringPage(p *agouti.Page) (*monitoringPage, error) {
	if err := p.Navigate(monitoringPageURL); err != nil {
		return nil, err
	}
	time.Sleep(10 * time.Second)
	return &monitoringPage{page: p}, nil
}

//
// Page components
//
func (p monitoringPage) settingToggleButton() *agouti.Selection {
	return p.page.FindByXPath("//*[@id=\"_0rif_sd-dashboard-toolbar\"]/button[2]")
}

func (p monitoringPage) oneColumnButton() *agouti.Selection {
	return p.page.FindByXPath("//*[@id=\"_0rif_mat-menu-panel-1\"]/div/button[6]")
}

func (p monitoringPage) maleThreeDotsToggleButton() *agouti.Selection {
	return p.page.FindByXPath("//*[@id=\"main\"]/div/div/central-page-area/div/div/pcc-content-viewport/div/div/pangolin-home/cfc-router-outlet/div/sd-dashboard-page/div/div/div/sd-dashboard-root/sd-grid/div/div[1]/sd-chart-card/mat-card/mat-card-header/sd-chart-header/div/div/sd-icon")
}

func (p monitoringPage) malePNGDownloadButton() *agouti.Selection {
	return p.page.FindByXPath("//*[@id=\"_0rif_mat-menu-panel-3\"]/div/div/button[4]")
}

func (p monitoringPage) femaleThreeDotsToggleButton() *agouti.Selection {
	return p.page.FindByXPath("//*[@id=\"main\"]/div/div/central-page-area/div/div/pcc-content-viewport/div/div/pangolin-home/cfc-router-outlet/div/sd-dashboard-page/div/div/div/sd-dashboard-root/sd-grid/div/div[2]/sd-chart-card/mat-card/mat-card-header/sd-chart-header/div/div/sd-icon")
}

func (p monitoringPage) femalePNGDownloadButton() *agouti.Selection {
	return p.page.FindByXPath("//*[@id=\"_0rif_mat-menu-panel-4\"]/div/div/button[4]")
}

func (p monitoringPage) oneDayButton() *agouti.Selection {
	return p.page.FindByXPath("//*[@id=\"_0rif_sd-dashboard-toolbar\"]/div[2]/button[3]")
}

func (p monitoringPage) oneWeekButton() *agouti.Selection {
	return p.page.FindByXPath("//*[@id=\"_0rif_sd-dashboard-toolbar\"]/div[2]/button[4]")
}

func (p monitoringPage) oneMonthButton() *agouti.Selection {
	return p.page.FindByXPath("//*[@id=\"_0rif_sd-dashboard-toolbar\"]/div[2]/button[5]")
}

//
// Page actions
//
func (p monitoringPage) download(rt ReportType) error {
	if err := p.settingToggleButton().Click(); err != nil {
		return fmt.Errorf("failed to click setting toggle button: %v", err)
	}
	time.Sleep(5 * time.Second)
	if err := p.oneColumnButton().Click(); err != nil {
		return fmt.Errorf("failed to click one column button: %v", err)
	}
	time.Sleep(5 * time.Second)

	switch rt {
	case DailyReport:
		if err := p.oneDayButton().Click(); err != nil {
			return fmt.Errorf("failed to click one day button: %v", err)
		}
	case WeeklyReport:
		if err := p.oneWeekButton().Click(); err != nil {
			return fmt.Errorf("failed to click one week button: %v", err)
		}
	case MonthlyReport:
		if err := p.oneMonthButton().Click(); err != nil {
			return fmt.Errorf("failed to click one month button: %v", err)
		}
	}

	if err := p.maleThreeDotsToggleButton().Click(); err != nil {
		return fmt.Errorf("failed to click male 3 dots toggle button: %v", err)
	}
	time.Sleep(5 * time.Second)
	if err := p.malePNGDownloadButton().Click(); err != nil {
		return fmt.Errorf("failed to click male png download button: %v", err)
	}
	time.Sleep(7 * time.Second)

	if err := p.femaleThreeDotsToggleButton().Click(); err != nil {
		return fmt.Errorf("failed to click female 3 dots toggle button: %v", err)
	}
	time.Sleep(5 * time.Second)
	if err := p.femalePNGDownloadButton().Click(); err != nil {
		return fmt.Errorf("failed to click female png download button: %v", err)
	}
	time.Sleep(7 * time.Second)

	return nil
}
