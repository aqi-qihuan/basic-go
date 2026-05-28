/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // HOK 营地风格 - 暗黑沉浸 + 金色荣耀
        brand: {
          gold: '#F0C060',
          'gold-light': '#F5D88A',
          'gold-dark': '#C8982A',
          orange: '#FF8C00',
          'orange-light': '#FFB347',
        },
        surface: {
          deep: '#0B0D17',
          card: '#131520',
          elevated: '#1A1D2B',
          hover: '#222536',
        },
        text: {
          primary: '#E8E0D0',
          bright: '#F5F0E8',
          secondary: '#9C9688',
          muted: '#6B6558',
          gold: '#F0C060',
        },
        tech: {
          blue: '#3B82F6',
          cyan: '#06B6D4',
          purple: '#8B5CF6',
        },
        rank: {
          gold: '#FFD700',
          silver: '#C0C0C0',
          bronze: '#CD7F32',
        },
        // 保留 primary 用于 Ant Design 兼容
        primary: {
          50: '#FFF7ED',
          100: '#FFEDD5',
          200: '#FED7AA',
          300: '#FDBA74',
          400: '#FB923C',
          500: '#F0C060',
          600: '#F0C060',
          700: '#C8982A',
          800: '#9A6E1A',
          900: '#7C5510',
        },
        success: '#22C55E',
        danger: '#EF4444',
        warning: '#F97316',
      },
      fontFamily: {
        sans: ['Inter', 'Noto Sans SC', 'PingFang SC', 'Microsoft YaHei', 'sans-serif'],
        heading: ['Inter', 'Noto Sans SC', 'PingFang SC', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
      borderRadius: {
        '4xl': '2rem',
      },
      boxShadow: {
        'glass': '0 4px 24px rgba(0, 0, 0, 0.4)',
        'glass-hover': '0 12px 40px rgba(0, 0, 0, 0.5), 0 0 20px rgba(240, 192, 96, 0.1)',
        'gold-glow': '0 0 30px rgba(240, 192, 96, 0.15)',
        'gold-btn': '0 4px 16px rgba(240, 192, 96, 0.3)',
        'gold-btn-hover': '0 8px 24px rgba(240, 192, 96, 0.4)',
        'rank-1': '0 0 20px rgba(255, 215, 0, 0.4)',
        'rank-2': '0 0 15px rgba(192, 192, 192, 0.3)',
        'rank-3': '0 0 15px rgba(205, 127, 50, 0.3)',
      },
      backdropBlur: {
        'glass': '24px',
      },
      animation: {
        'gold-pulse': 'goldPulse 2s ease-in-out infinite',
        'fade-in-up': 'fadeInUp 0.5s ease-out',
        'shimmer': 'shimmer 1.5s infinite',
      },
      keyframes: {
        goldPulse: {
          '0%, 100%': { boxShadow: '0 0 20px rgba(255, 215, 0, 0.4)' },
          '50%': { boxShadow: '0 0 40px rgba(255, 215, 0, 0.2)' },
        },
        fadeInUp: {
          from: { opacity: '0', transform: 'translateY(20px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
        shimmer: {
          '0%': { backgroundPosition: '-1000px 0' },
          '100%': { backgroundPosition: '1000px 0' },
        },
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
}
