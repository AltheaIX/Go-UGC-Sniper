package main

import (
	"testing"
)

func TestGetAccountDetails(t *testing.T) {
	GetAccountDetails("_|WARNING:-DO-NOT-SHARE-THIS.--Sharing-this-will-allow-someone-to-log-in-as-you-and-to-steal-your-ROBUX-and-items.|_4A1604EC09ACC1529A175808A10D2A9586BB0ED4DC9191E11515C44B56A037C15EEFF6ADD523C6F45ADDF9A9A3D09EB7D630EAEBA0D89192EC05B5143B3B93DFBF763CDA6B014019317B728A691F73016B1949F78916AC6581B39BC748D69B6529811FB814E13003DAF6FBF1DA8EE089A4D54D435798A929A16702E6267AB376625BB0741B6F3CCE7842D98EF822D03B2465D0C797E47FC407841C847B4ED80898A1CDD4797955109FED0501776660CC80FF7C8F1EFE6C7DE8A1252CA865189097B2C8B5ED4B9F77ACAA1B1EE737EC3191BCC4432E6F46FDB58753740E44CFC2BE7AEB294D758047DFF82D4C5ABCDC0D2FCB7C82856124E4435DB803734682E6B35465871D4D455CDBC2991160D94F1EA3EC6F41611977854BA7CADE34EF24769C06C48ED5CBE222A931394E32BE33014D34368DCBB81FB63112C0156C90D86E2FFEED45F6487019BF2C8A1B0FA68CD58FE348102747DC00080A7C6561423B800A9390E66AC2A33FC318EAC0707AAF4194A2FB6E")
	t.Log(accountId)
}

func TestLoadConfig(t *testing.T) {
	LoadConfig()
}
