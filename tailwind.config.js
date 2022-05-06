module.exports = {
    purge: {
        enabled: true,
        content: [
            './content/**/*.md',
            './layouts/**/*.html',
        ],
    },
    theme: {
        container: {
            center: true,
            padding: {
                default: '1rem',
                sm: '2rem',
                lg: '4rem',
                xl: '6rem',
                '2xl': '6rem',
            },
        },
        extend: {
            backgroundColor: {
                primary: "var(--color-bg-primary)",
                secondary: "var(--color-bg-secondary)",
            },
            textColor: {
                accent: "var(--color-text-accent)",
                primary: "var(--color-text-primary)",
                secondary: "var(--color-text-secondary)",
            },
            stroke: {
                current: "var(--color-text-primary)",
            }
        },
    },
    variants: {},
    corePlugins: {},
    plugins: [],
}