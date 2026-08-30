/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        bg: '#050505',
        panel: '#121212',
        borderRed: '#2a0808',
        primaryRed: '#ff1e1e',
        darkRed: '#8b0000',
      },
    },
  },
  plugins: [],
};
