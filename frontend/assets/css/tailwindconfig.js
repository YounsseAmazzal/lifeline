        tailwind.config = {
            theme: {
                extend: {
                    fontFamily: {
                        sans: ['"Plus Jakarta Sans"', 'sans-serif'],
                        arabic: ['"Cairo"', 'sans-serif'],
                    },
                    colors: {
                        morocco: {
                            red: '#C1272D',     // Ahmar
                            green: '#006233',   // Akhdar
                            gold: '#D4AF37',    // Dahbi
                            dark: '#0f172a',
                        }
                    },
                    backgroundImage: {
                        'zellige': "url('https://www.transparenttextures.com/patterns/arabesque.png')",
                        'royal-gradient': "linear-gradient(135deg, rgba(193, 39, 45, 0.05) 0%, rgba(0, 98, 51, 0.05) 100%)",
                    },
                    animation: {
                        'scroll': 'scroll 40s linear infinite', /* Speed of logos */
                        'float': 'float 6s ease-in-out infinite',
                        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
                    },
                    keyframes: {
                        scroll: {
                            '0%': { transform: 'translateX(0)' },
                            '100%': { transform: 'translateX(50%)' }, /* RTL moves right, use positive 50% */
                        },
                        float: {
                            '0%, 100%': { transform: 'translateY(0)' },
                            '50%': { transform: 'translateY(-20px)' },
                        }
                    }
                }
            }
        }