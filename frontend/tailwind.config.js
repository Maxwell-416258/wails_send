/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
    function ({ addUtilities }) {
      addUtilities({
        '.input-glass': {
          'width': '24rem', // w-96
          'padding': '0.5rem 1rem', // px-4 py-2
          'borderRadius': '0.75rem', // rounded-xl
          'backgroundColor': 'rgba(255, 255, 255, 0.05)', // bg-white/10
          'color': 'white',
          'border': '1px solid rgba(255, 255, 255, 0.2)', // border-white/20
          'backdropFilter': 'blur(12px)', // backdrop-blur-md
        },
        '.input-glass:focus': {
          'outline': 'none',
          'boxShadow': '0 0 0 2px rgba(192, 132, 252, 0.5)', // 类似 focus:ring-purple-400
        },
      })
    }
  ],
}


