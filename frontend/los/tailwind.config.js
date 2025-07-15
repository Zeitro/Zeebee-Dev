module.exports = {
  content: [
    "./src/**/*.{js,ts,jsx,tsx}", // 根据你的项目结构调整
  ],
 theme: {
    extend: {
      colors: {
        background: 'var(--background)',
        foreground: 'var(--foreground)',
        base: {
          primary: 'var(--base-primary)',
          success: 'var(--base-success)',
          warning: 'var(--base-warning)',
          error: 'var(--base-error)',
          disabled: 'var(--base-disabled)',
        },
        black: {
          light: 'var(--black-light)', // 浅黑色
          medium: 'var(--black-medium)', // 中黑色
          dark: 'var(--black-dark)', // 深黑色
        },
        gray: {// 灰色调
          light: 'var(--gray-light)', // 浅灰色
          medium: 'var(--gray-medium)',// 中灰色
          dark: 'var(--gray-dark)', // 深灰色
        },
        blue: {
          light: 'var(--blue-light)',
          medium: 'var(--blue-medium)',
          dark: 'var(--blue-dark)',
        },
      },
      fontFamily: {
        sans: ['var(--font-inter)', 'Arial', 'sans-serif'],
        playfair: ['var(--font-playfair)', 'serif'],
        lexend: ['var(--font-lexend)', 'sans-serif'],
        openSans: ['var(--font-open-sans)', 'sans-serif'],
        roboto: ['var(--font-roboto)'],
      },
    },
  },
  plugins: [],
};