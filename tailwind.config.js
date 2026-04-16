/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: ["./web/templates/**/*.templ", "./web/templates/**/*_templ.go"],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        display: ['Phudu', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      colors: {
        ivory: '#FAF8F3',
        surface: '#FFFFFF',
        ink: '#1A1918',
        forest: {
          DEFAULT: '#1F3D2E',
          50: '#E8EDE9',
          100: '#C4D1C8',
          600: '#2D5C44',
          900: '#122419',
        },
        ember: '#D94F30',
        flame: '#E8623A',
        muted: '#6B665E',
        hairline: '#E8E4DB',
        cream: '#F2EEE4',
        sage: '#4A7C59',
        rust: '#A64545',
        // Dark palette
        'dark-bg': '#0F1117',
        'dark-surface': '#1A1D26',
        'dark-text': '#E8E6E3',
        'dark-border': '#2D3039',
        'dark-cream': '#1F2229',
        'dark-muted': '#8B8880',
      }
    }
  },
  plugins: [],
}
