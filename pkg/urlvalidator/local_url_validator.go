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
	// CommonJSLibraryFileRegex is a regex to match common js library files.
	CommonJSLibraryFileRegex         = `(?i)(?:amplify|quantserve|slideshow|jquery|modernizr|polyfill|vendor|modules|gtm|underscore?|tween|retina|selectivizr|cufon|angular|swf|sha1|freestyle|bootstrap|d3|backbone|videojs|google[-_]analytics|material|redux|knockout|datepicker|datetimepicker|ember|react|ng|fusion|analytics|libs?|vendors?|node[-_]modules|lodash|moment|chart|highcharts|raphael|prototype|mootools|dojo|ext|yui|web[-_]?components|polymer|vue|svelte|next|nuxt|gatsby|express|koa|hapi|socket[-_.]?io|axios|superagent|request|bluebird|rxjs|ramda|immutable|flux|redux[-_]saga|mobx|relay|apollo|graphql|three|phaser|pixi|babylon|cannon|hammer|howler|gsap|velocity|mo[-_.]?js|popper|shepherd|prism|highlight|markdown[-_]?it|codemirror|ace[-_]?editor|tinymce|ckeditor|quill|simplemde|monaco[-_]?editor|pdf[-_.]?js|jspdf|fabric|paper|konva|p5|processing|matter[-_.]?js|box2d|planck|chart[-_.]?js|plotly|echarts|d3[-_.]?force|sigma|c3|nvd3|amcharts|vis[-_.]?js|dagre[-_.]?d3|cytoscape|leaflet|openlayers|ol3|mapbox|cesium|turf|moment[-_.]?timezone|luxon|dayjs|date[-_.]?fns|date[-_.]?io|flatpickr|pikaday|fullcalendar|draggable|interact|sortable|dragula|dropzone|filepond|uppy|fine[-_.]?uploader|plyr|mediaelement|flowplayer|jwplayer|video[-_.]?js|mediaelement[-_.]?js|dash[-_.]?js|hls[-_.]?js|videojs|wavesurfer|soundmanager|amplitude|pizzicato|tone|adroll|doubleclick|facebook-pixel|ga-audiences|googlesyndication|adsbygoogle|gpt|amazon-adsystem|criteo|taboola|outbrain|bidswitch|bidswitch.net|spotxchange|yahoo|media.net|contextweb|openx|pubmatic|rubiconproject|indexexchange|appnexus|liveintent|triplelift|verizonmedia|synacor|sonobi|yieldmo|gumgum|smartadserver|mopub|pubnative|inmobi|chartboost|tapjoy|admob|unityads|vungle|flurry|matomy|altitude|dataxu|thetradedesk|exponential|zypmedia|quantcast|mediamath|bidswitch|mgid|revcontent|powerlinks|rhythmone|airpush|smaato|adcolony|mopub|leadbolt|mobfox|nativo|revjet|smartyads|avocarrot|epom|imobile|supersonicads|loopme|applovin|pandora|mytarget|bidvertiser|chitika|popads|propellerads|buysellads|adhit|hilltopads|plugrush|popcash|popunder|revenuehits|trafficjunky|trafficfactory|zero-|smartoasis)(?:[-._][\w\d]*)*\.js$`
	commonJSLibraryFileRegexCompiled = regexp.MustCompile(CommonJSLibraryFileRegex)
)

func ValidateLocalUrl(config *config.Config, parsedUrl *url.URL, domain string, registeredDomain string) (*url.URL, error) {

	urlExt := filepath.Ext(parsedUrl.Path)

	for _, allwdExt := range config.AllowedExtensions {

		if urlExt == allwdExt {

			switch allwdExt {

			case "", ".htm", ".html":
				return handleSubpages(config, parsedUrl)
			case ".js":
				return handleJavascript(parsedUrl)
			default:
				return parsedUrl, nil

			}

		}
	}

	return nil, fmt.Errorf("Url has unallowed extension")
}

func handleSubpages(config *config.Config, parsedUrl *url.URL) (*url.URL, error) {

	urlPathDepth := getPathDepth(parsedUrl.Path)

	if urlPathDepth <= config.MaxPathDepth {
		return parsedUrl, nil

	} else if hasSensitivePattern(config, parsedUrl.Path) {
		return parsedUrl, nil
	}

	return nil, fmt.Errorf("Too deep without sensitive patterns")
}

func getPathDepth(urlPath string) int {

	trimmedPath := strings.Trim(urlPath, `/`)
	if trimmedPath == "" {
		return 0
	}

	pathDepth := len(strings.Split(trimmedPath, `/`))

	return pathDepth

}

func hasSensitivePattern(config *config.Config, urlPath string) bool {

	for _, pttn := range config.SensitivePatterns {

		if strings.Contains(urlPath, pttn) {
			return true
		}

	}

	return false
}

func handleJavascript(parsedUrl *url.URL) (*url.URL, error) {

	if IsPathCommonJSLibraryFile(parsedUrl.Path) {
		return nil, fmt.Errorf("Common .js library")

	} else {
		return parsedUrl, nil
	}
}

// IsPathCommonJSLibraryFile checks if a given path is a common js library file.
// Taken from
// https://github.com/projectdiscovery/katana/blob/462965a11211371b9f3fd6b2ef6c51ee9913de90/pkg/utils/jsluice.go#L14
func IsPathCommonJSLibraryFile(urlPath string) bool {
	return commonJSLibraryFileRegexCompiled.MatchString(urlPath)
}
