/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: '#A7FF83', // Neon lime
        'primary-dark': '#76C958',
        background: '#000000',
        surface: '#121212',
        'surface-highlight': '#1E1E1E',
      },
    },
  },
  plugins: [require('daisyui')],
  daisyui: {
    themes: [
      {
        dark: {
          ...require('daisyui/src/theming/themes')['dark'],
          primary: '#A7FF83',
          'primary-content': '#000000',
          'base-100': '#000000',
          'base-200': '#121212',
          'base-300': '#1E1E1E',
        },
      },
    ],
    darkTheme: 'dark',
  },
}

