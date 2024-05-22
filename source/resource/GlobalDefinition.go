package resource

import "fyne.io/fyne/v2"

//Put all the text data here so we can configure the language setting at runtime.
const KGoldenRatio = 1.618033988749
const KGoldenRatioPercentage = 1 / (1 + KGoldenRatio)

var KWindowSize fyne.Size = fyne.NewSize(600, 650)
var KAlbumCardSize fyne.Size = fyne.NewSize(140, 180)
var KAlbumCoverSize fyne.Size = fyne.NewSize(160, 160)

//Program
const KWindowTitle = "MeowyPlayer"

//Generic Text
const KUploadText = "Upload"
const KCancelText = "Cancel"
const KCreateText = "Create"
const KDeleteText = "Delete"
const KSaveText = "Save"
const KEditText = "Edit"
const KHomeText = "Home"
const KCollectionText = "Collection"
const KAccountText = "Account"
const KSettingText = "Setting"
const KReturnText = "Return"

const KSortMenuText = "Sort By"
const KMostRecentText = "Most Recent"
const KAlphabeticalText = "Alphabetical"
const KCreateAlbumText = "Create Album"
const KEditAlbumText = "Edit Album"
const KEnterTitleHint = "Enter title:"
const KDeleteConfirmationText = "Delete Confirmation"
const KDeleteAlbumTextTemplate = "Delete the album \"%v\" ?"
const KDeleteMusicTextTemplate = "Delete the music \"%v\" ?"
const KAlbumTipTextTemplate = "Songs: %v\n\nLast Edit: %v"
