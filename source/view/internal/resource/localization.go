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

func WindowTitle() string            { return lang.L("MeowyPlayer") }
func DownloadText() string           { return lang.L("Download") }
func UploadText() string             { return lang.L("Upload") }
func CancelText() string             { return lang.L("Cancel") }
func CreateText() string             { return lang.L("Create") }
func DeleteText() string             { return lang.L("Delete") }
func SaveText() string               { return lang.L("Save") }
func EditText() string               { return lang.L("Edit") }
func HomeText() string               { return lang.L("Home") }
func CollectionsText() string        { return lang.L("Collections") }
func AccountText() string            { return lang.L("Account") }
func SettingsText() string           { return lang.L("Settings") }
func BackText() string               { return lang.L("Back") }
func SearchingText() string          { return lang.L("Searching") }
func AttemptsText() string           { return lang.L("Attempts") }
func SortMenuText() string           { return lang.L("Sort By") }
func MostRecentText() string         { return lang.L("Most Recent") }
func AlphabeticalText() string       { return lang.L("Alphabetical") }
func CreateAlbumText() string        { return lang.L("Create Album") }
func EditAlbumText() string          { return lang.L("Edit Album") }
func EnterTitleHint() string         { return lang.L("Enter title:") }
func DeleteConfirmationText() string { return lang.L("Delete Confirmation") }
func SelectAlbumText() string        { return lang.L("(Select an album)") }
func RandomText() string             { return lang.L("Random") }
func OrderedText() string            { return lang.L("Ordered") }
func RepeatText() string             { return lang.L("Repeat") }
func UsernameText() string           { return lang.L("Username") }
func PasswordText() string           { return lang.L("Password") }
func LoginText() string              { return lang.L("Login") }
func LogoutText() string             { return lang.L("Logout") }
func RegisterText() string           { return lang.L("Register") }
func LoginToContinueText() string    { return lang.L("Login to continue") }
func MigrateToRemoteText() string    { return lang.L("Migrate to remote") }
func MigrateAlbumsText() string      { return lang.L("Migrate Albums") }
func MigrateConfirmationText() string {
	return lang.L("Do you want to upload all the local albums to remote?")
}
func DeleteAlbumTextTemplate() string { return lang.L("Delete the album \"%v\" ?") }
func DeleteMusicTextTemplate() string { return lang.L("Delete the music \"%v\" ?") }
func AlbumTipTextTemplate() string    { return lang.L("Songs: %v\n\nLast Edit: %v") }
