/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/templates/**/*.tmpl"
  ],
  theme: {
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
  plugins: [],
}

