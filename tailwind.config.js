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
        // CSS-variable-based colors: auto-switch light↔dark
        ivory: 'var(--bg)',
        surface: 'var(--surface)',
        ink: 'var(--ink)',
        muted: 'var(--muted)',
        hairline: 'var(--hairline)',
        cream: 'var(--cream)',

        // Accent colors: use dark-aware values
        forest: {
          DEFAULT: 'var(--forest)',
          50: 'var(--forest-50)',
          100: '#C4D1C8',
          600: 'var(--forest-600)',
          900: '#122419',
        },
        ember: 'var(--ember)',
        flame: 'var(--flame)',
        sage: 'var(--sage)',
        rust: 'var(--rust)',

        // Explicit dark palette (for dark: prefix usage)
        'dark-bg': '#0D0F14',
        'dark-surface': '#161920',
        'dark-text': '#F2F0ED',
        'dark-border': '#2A2E36',
        'dark-cream': '#1C1F27',
        'dark-muted': '#B0AAA0',
      }
    }
  },
  plugins: [],
}
