import forms from '@tailwindcss/forms'

/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx,js,jsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: '#f97316',
        'primary-strong': '#c5561a',
        accent: '#facc15',
        'accent-soft': '#fef3c7',
        'background-light': '#faf7f2',
        'background-muted': '#ede9e4',
        'background-dark': '#0f172a',
        'surface-dark': '#162235',
      },
      fontFamily: {
        display: ['Inter', 'sans-serif'],
      },
      borderRadius: {
        lg: '0.5rem',
        xl: '0.75rem',
      },
      boxShadow: {
        card: '0 10px 30px -12px rgba(20, 50, 80, 0.25)',
      },
    },
  },
  plugins: [forms],
}
