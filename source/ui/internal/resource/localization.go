package resource

import (
	"embed"

	"fyne.io/fyne/v2/lang"
)

//go:embed translation
var translations embed.FS

func RegisterTranslation() {
	lang.AddTranslationsFS(translations, "translation")
}

func AccountText() string                       { return lang.L("Account") }
func AlbumTipTextTemplate() string              { return lang.L("Songs: %v\n\nLast Edit: %v") }
func AlphabeticalText() string                  { return lang.L("Alphabetical") }
func AttemptsText() string                      { return lang.L("Attempts") }
func BackText() string                          { return lang.L("Back") }
func CancelText() string                        { return lang.L("Cancel") }
func CollectionsText() string                   { return lang.L("Collections") }
func CreateAlbumText() string                   { return lang.L("Create Album") }
func CreateText() string                        { return lang.L("Create") }
func DeleteAlbumTextTemplate() string           { return lang.L("Delete the album \"%v\" ?") }
func DeleteConfirmationText() string            { return lang.L("Delete Confirmation") }
func DeleteMusicTextTemplate() string           { return lang.L("Delete the music \"%v\" ?") }
func DeleteText() string                        { return lang.L("Delete") }
func DownloadText() string                      { return lang.L("Download") }
func EditAlbumText() string                     { return lang.L("Edit Album") }
func EditText() string                          { return lang.L("Edit") }
func EnterTitleHint() string                    { return lang.L("Enter title:") }
func HomeText() string                          { return lang.L("Home") }
func LoginText() string                         { return lang.L("Login") }
func LoginToContinueText() string               { return lang.L("Login to continue") }
func LogoutText() string                        { return lang.L("Logout") }
func MigrateConfirmationText() string           { return lang.L("Confirm to migrate the albums?") }
func UploadLocalAlbumsToTheAccountText() string { return lang.L("Upload local albums to the account") }
func BackupAlbumsToLocalText() string           { return lang.L("Backup albums to local") }
func MostRecentText() string                    { return lang.L("Most Recent") }
func OrderedText() string                       { return lang.L("Ordered") }
func PasswordText() string                      { return lang.L("Password") }
func RandomText() string                        { return lang.L("Random") }
func RegisterText() string                      { return lang.L("Register") }
func RepeatText() string                        { return lang.L("Repeat") }
func SaveText() string                          { return lang.L("Save") }
func SearchingText() string                     { return lang.L("Searching") }
func SelectAlbumText() string                   { return lang.L("Select an album") }
func SettingsText() string                      { return lang.L("Settings") }
func SortMenuText() string                      { return lang.L("Sort By") }
func UploadText() string                        { return lang.L("Upload") }
func UsernameText() string                      { return lang.L("Username") }
func WindowTitle() string                       { return lang.L("MeowyPlayer") }
