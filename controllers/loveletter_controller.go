package controllers

import (
	"github.com/astaxie/beego"
	"github.com/ulricqin/beego-blog/models/loveletter"
)

type LoveLetterController struct {
	beego.Controller
}

func (this *LoveLetterController) Get() {
	this.Data["NrLetter"] = loveletter.GetNrLetter()
	this.TplNames = "editletter.html"

}

func (this *LoveLetterController) Post() {
	content := this.GetString("letter")

	beego.Info("redirect to edit ")
	loveletter.Add(content)
	this.Redirect("/editletter", 302)
}
