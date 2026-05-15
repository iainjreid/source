/** @type {import('tailwindcss').Config} */
module.exports = {
  experimental: {
    optimizeUniversalDefaults: true,
  },
  content: [
    "./web/templates/**/*.tmpl"
  ],
  theme: {
    container: {
      center: true,
    },
    fontFamily: {
      mono: [
        '"CommitMono", monospace',
        {
          fontFeatureSettings: '"cv10", "ss01"',
          fontVariationSettings: '"opsz" 32',
        },
      ],
    },
    fontVariationSettings: {"weight":700,"italic":false,"alternates":{"cv01":false,"cv02":false,"cv03":false,"cv04":false,"cv05":false,"cv06":false,"cv07":false,"cv08":false,"cv09":false,"cv10":true,"cv11":false},"features":{"ss01":false,"ss02":false,"ss03":true,"ss04":true,"ss05":true},"letterSpacing":0,"lineHeight":1}
  },
  safelist: [
    "group/1", "group-open/1:flex", "group-open/1:hidden",
    "group/2", "group-open/2:flex", "group-open/2:hidden",
    "group/3", "group-open/3:flex", "group-open/3:hidden",
    "group/4", "group-open/4:flex", "group-open/4:hidden",
    "group/5", "group-open/5:flex", "group-open/5:hidden",
    "group/6", "group-open/6:flex", "group-open/6:hidden",
    "group/7", "group-open/7:flex", "group-open/7:hidden",
    "group/8", "group-open/8:flex", "group-open/8:hidden",
    "group/9", "group-open/9:flex", "group-open/9:hidden",
    {
      pattern: /pl-(1|2|3|4|5|6|7|8|9)/
    }
  ],
}

