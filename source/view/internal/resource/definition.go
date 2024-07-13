package resource

import (
	"fyne.io/fyne/v2"
)

// Put all the data here so we can configure the language setting at runtime.
const KGoldenRatio = 1.618033988749
const KGoldenRatioPercentage = 1 / (1 + KGoldenRatio)

var KWindowSize fyne.Size = fyne.NewSize(600, 650)
var KAlbumCardSize fyne.Size = fyne.NewSize(140, 180)
var KAlbumCoverSize fyne.Size = fyne.NewSize(160, 160)
var KAlbumPreviewSize fyne.Size = fyne.NewSize(140, 140)
var KThumbnailSize fyne.Size = fyne.NewSize(271, 152)

const KSearchAttempts = 5

// Program Text
const KWindowTitle = "MeowyPlayer"
const KDownloadText = "Download"
const KUploadText = "Upload"
const KQueueText = "Queue"
const KCancelText = "Cancel"
const KCreateText = "Create"
const KDeleteText = "Delete"
const KSaveText = "Save"
const KEditText = "Edit"
const KHomeText = "Home"
const KCollectionText = "Collection"
const KAccountText = "Account"
const KSettingText = "Setting"
const KBackText = "Back"
const KSearchingText = "Searching"
const KAttemptsText = "Attempts"
const KSortMenuText = "Sort By"
const KMostRecentText = "Most Recent"
const KAlphabeticalText = "Alphabetical"
const KCreateAlbumText = "Create Album"
const KEditAlbumText = "Edit Album"
const KEnterTitleHint = "Enter title:"
const KDeleteConfirmationText = "Delete Confirmation"
const KSelectAlbumText = "(Select an album)"
const KRandomText = "Random"
const KOrderedText = "Ordered"
const KRepeatText = "Repeat"

const KDeleteAlbumTextTemplate = "Delete the album \"%v\" ?"
const KDeleteMusicTextTemplate = "Delete the music \"%v\" ?"
const KAlbumTipTextTemplate = "Songs: %v\n\nLast Edit: %v"
