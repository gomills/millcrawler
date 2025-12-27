package urlvalidator

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gomills/gofocusedcrawler/internal/config"
)

var (
	// CommonJSLibraryFileRegex is a regex to match common js library files and skip them.
	// Source: https://github.com/projectdiscovery/katana
	CommonJSLibraryFileRegex = `(?i)(?:amplify|quantserve|slideshow|jquery|modernizr|polyfill|vendor|modules|gtm|underscore?|tween|retina|selectivizr|cufon|angular|swf|sha1|freestyle|bootstrap|d3|backbone|videojs|google[-_]analytics|material|redux|knockout|datepicker|datetimepicker|ember|react|ng|fusion|analytics|libs?|vendors?|node[-_]modules|lodash|moment|chart|highcharts|raphael|prototype|mootools|dojo|ext|yui|web[-_]?components|polymer|vue|svelte|next|nuxt|gatsby|express|koa|hapi|socket[-_.]?io|axios|superagent|request|bluebird|rxjs|ramda|immutable|flux|redux[-_]saga|mobx|relay|apollo|graphql|three|phaser|pixi|babylon|cannon|hammer|howler|gsap|velocity|mo[-_.]?js|popper|shepherd|prism|highlight|markdown[-_]?it|codemirror|ace[-_]?editor|tinymce|ckeditor|quill|simplemde|monaco[-_]?editor|pdf[-_.]?js|jspdf|fabric|paper|konva|p5|processing|matter[-_.]?js|box2d|planck|chart[-_.]?js|plotly|echarts|d3[-_.]?force|sigma|c3|nvd3|amcharts|vis[-_.]?js|dagre[-_.]?d3|cytoscape|leaflet|openlayers|ol3|mapbox|cesium|turf|moment[-_.]?timezone|luxon|dayjs|date[-_.]?fns|date[-_.]?io|flatpickr|pikaday|fullcalendar|draggable|interact|sortable|dragula|dropzone|filepond|uppy|fine[-_.]?uploader|plyr|mediaelement|flowplayer|jwplayer|video[-_.]?js|mediaelement[-_.]?js|dash[-_.]?js|hls[-_.]?js|videojs|wavesurfer|soundmanager|amplitude|pizzicato|tone|adroll|doubleclick|facebook-pixel|ga-audiences|googlesyndication|adsbygoogle|gpt|amazon-adsystem|criteo|taboola|outbrain|bidswitch|bidswitch.net|spotxchange|yahoo|media.net|contextweb|openx|pubmatic|rubiconproject|indexexchange|appnexus|liveintent|triplelift|verizonmedia|synacor|sonobi|yieldmo|gumgum|smartadserver|mopub|pubnative|inmobi|chartboost|tapjoy|admob|unityads|vungle|flurry|matomy|altitude|dataxu|thetradedesk|exponential|zypmedia|quantcast|mediamath|bidswitch|mgid|revcontent|powerlinks|rhythmone|airpush|smaato|adcolony|mopub|leadbolt|mobfox|nativo|revjet|smartyads|avocarrot|epom|imobile|supersonicads|loopme|applovin|pandora|mytarget|bidvertiser|chitika|popads|propellerads|buysellads|adhit|hilltopads|plugrush|popcash|popunder|revenuehits|trafficjunky|trafficfactory|zero-|smartoasis)(?:[-._][\w\d]*)*\.js$`
	commonJSLibFilesCompReg  = regexp.MustCompile(CommonJSLibraryFileRegex)
)

// validateLocalUrl takes a local url and checks if it passes heuristics to be valid
func validateLocalUrl(config *config.Config, parsedUrl *url.URL, domain string, registeredDomain string) (*url.URL, error) {

	if parsedUrl == nil {
		return nil, fmt.Errorf("parsed_url_is_nil")
	}

	// get its extension and perform the check according to it (scripts, subpages or others .txt etc)
	urlExt := filepath.Ext(parsedUrl.Path)

	for _, allwdExt := range config.AllowedExtensions {

		if urlExt == allwdExt {

			switch allwdExt {

			// these extensions belong to subpages, fetching .html pages
			case "", ".htm", ".html":
				return validateSubpage(config, parsedUrl)

			// javascript files
			case ".js":
				return validateJs(parsedUrl)

			// here fall .yaml, .txt, .map, etc
			default:
				return parsedUrl, nil

			}

		}
	}

	// if it wasn't found, it wasn't in the allowed extensions list (.png, .svg, etc)
	return nil, fmt.Errorf("url_no_allowed_extension")
}

// validateSubpage validates subpages, which are .html files.
func validateSubpage(config *config.Config, parsedUrl *url.URL) (*url.URL, error) {

	if parsedUrl == nil {
		return nil, fmt.Errorf("parsed_url_is_nil")
	}

	// 1. get its path depth
	urlPathDepth := getPathDepth(parsedUrl.Path)

	// 2.0 if it's in the allowed range it's valid instantly
	if urlPathDepth <= config.MaxPathDepth {
		return parsedUrl, nil

		// 2.1 otherwise only valid if it has a sensitive pattern
	} else if hasSensitivePattern(config, parsedUrl.Path) {
		return parsedUrl, nil
	}

	return nil, fmt.Errorf("Too deep without sensitive patterns")
}

// getPathDepth calculates path depth of a subpage (e.g: 2 for example.com/path1/path2)
func getPathDepth(urlPath string) int {

	// get rid of trailing slashes (e.g: path1/ -> path1)
	trimmedPath := strings.Trim(urlPath, `/`)
	if trimmedPath == "" {
		return 0
	}

	pathDepth := len(strings.Split(trimmedPath, `/`))

	return pathDepth

}

// hasSensitivePattern returns true if there's a sensitive pattern in the url's path
func hasSensitivePattern(config *config.Config, urlPath string) bool {

	for _, pttn := range config.SensitivePatterns {

		if strings.Contains(urlPath, pttn) {
			return true
		}

	}

	return false
}

// validateJs checks that the .js file is not from a common library to avoid resources waste. The value
// of this check is that we regex here the url path to avoid working with whole javascript files of these
// common libraries in the future
func validateJs(parsedUrl *url.URL) (*url.URL, error) {

	if parsedUrl == nil {
		return nil, fmt.Errorf("parsed_url_is_nil")
	}

	if IsPathCommonJSLibraryFile(parsedUrl.Path) {
		return nil, fmt.Errorf("common_js_library")

	} else {
		return parsedUrl, nil
	}
}

// IsPathCommonJSLibraryFile checks if a given path is a common js library file.
// Taken from
// https://github.com/projectdiscovery/katana/blob/462965a11211371b9f3fd6b2ef6c51ee9913de90/pkg/utils/jsluice.go#L14
func IsPathCommonJSLibraryFile(urlPath string) bool {
	return commonJSLibFilesCompReg.MatchString(urlPath)
}
