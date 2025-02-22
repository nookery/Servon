// 语言配置接口
export interface LanguageConfig {
    id: string
    name: string
    extensions: string[]
    filenames?: string[]
}

// 语言配置列表
export const languageConfigs: LanguageConfig[] = [
    // 常见文本文件
    { id: 'plaintext', name: '纯文本', extensions: ['txt'], filenames: ['.gitignore', '.env', '.editorconfig'] },
    { id: 'markdown', name: 'Markdown', extensions: ['md', 'markdown', 'mdown', 'mkdn', 'mkd'] },

    // 配置文件
    { id: 'json', name: 'JSON', extensions: ['json', 'jsonc', 'json5'], filenames: ['package.json', 'package-lock.json', 'composer.json'] },
    { id: 'yaml', name: 'YAML', extensions: ['yaml', 'yml'], filenames: ['docker-compose.yml', '.travis.yml'] },
    { id: 'toml', name: 'TOML', extensions: ['toml'] },
    { id: 'ini', name: 'INI', extensions: ['ini', 'conf', 'cfg', 'prefs'] },

    // Web 开发
    { id: 'javascript', name: 'JavaScript', extensions: ['js', 'jsx', 'mjs', 'cjs'] },
    { id: 'typescript', name: 'TypeScript', extensions: ['ts', 'tsx', 'mts', 'cts'] },
    { id: 'vue', name: 'Vue', extensions: ['vue'] },
    { id: 'html', name: 'HTML', extensions: ['html', 'htm', 'xhtml', 'shtml'] },
    { id: 'css', name: 'CSS', extensions: ['css'] },
    { id: 'scss', name: 'SCSS', extensions: ['scss', 'sass'] },
    { id: 'less', name: 'Less', extensions: ['less'] },
    { id: 'jsx', name: 'React JSX', extensions: ['jsx'] },
    { id: 'tsx', name: 'React TSX', extensions: ['tsx'] },

    // 后端语言
    { id: 'python', name: 'Python', extensions: ['py', 'pyw', 'pyx', 'pxd', 'pxi', 'pyc'] },
    { id: 'java', name: 'Java', extensions: ['java', 'jav'] },
    { id: 'go', name: 'Go', extensions: ['go'] },
    { id: 'rust', name: 'Rust', extensions: ['rs'] },
    { id: 'php', name: 'PHP', extensions: ['php', 'php4', 'php5', 'phtml', 'ctp'] },
    { id: 'ruby', name: 'Ruby', extensions: ['rb', 'rbw', 'rake', 'gemspec'] },
    { id: 'perl', name: 'Perl', extensions: ['pl', 'pm', 'pod', 't'] },
    { id: 'kotlin', name: 'Kotlin', extensions: ['kt', 'kts'] },
    { id: 'swift', name: 'Swift', extensions: ['swift'] },
    { id: 'csharp', name: 'C#', extensions: ['cs'] },
    { id: 'cpp', name: 'C++', extensions: ['cpp', 'cc', 'cxx', 'c++', 'hpp', 'hh', 'hxx', 'h++'] },
    { id: 'c', name: 'C', extensions: ['c', 'h'] },

    // 数据库
    { id: 'sql', name: 'SQL', extensions: ['sql'] },
    { id: 'mysql', name: 'MySQL', extensions: ['mysql'] },
    { id: 'pgsql', name: 'PostgreSQL', extensions: ['pgsql'] },

    // 标记语言
    { id: 'xml', name: 'XML', extensions: ['xml', 'xsd', 'xsl', 'xslt', 'svg'] },
    { id: 'latex', name: 'LaTeX', extensions: ['tex', 'cls', 'sty'] },

    // Shell 脚本
    { id: 'shell', name: 'Shell Script', extensions: ['sh', 'bash', 'zsh', 'fish'], filenames: ['.bashrc', '.zshrc'] },
    { id: 'powershell', name: 'PowerShell', extensions: ['ps1', 'psm1', 'psd1'] },
    { id: 'batch', name: 'Batch', extensions: ['bat', 'cmd'] },

    // 其他
    { id: 'dockerfile', name: 'Dockerfile', extensions: [], filenames: ['dockerfile', 'Dockerfile'] },
    { id: 'makefile', name: 'Makefile', extensions: ['mk'], filenames: ['makefile', 'Makefile'] },
    { id: 'groovy', name: 'Groovy', extensions: ['groovy', 'gvy', 'gy', 'gsh'] },
    { id: 'r', name: 'R', extensions: ['r', 'R'] },
    { id: 'dart', name: 'Dart', extensions: ['dart'] },
    { id: 'lua', name: 'Lua', extensions: ['lua'] },
    { id: 'scala', name: 'Scala', extensions: ['scala', 'sc'] },
    { id: 'haskell', name: 'Haskell', extensions: ['hs', 'lhs'] },
    { id: 'elixir', name: 'Elixir', extensions: ['ex', 'exs'] },
    { id: 'erlang', name: 'Erlang', extensions: ['erl', 'hrl'] },
    { id: 'clojure', name: 'Clojure', extensions: ['clj', 'cljs', 'cljc', 'edn'] },

    // 配置和依赖文件
    { id: 'properties', name: 'Properties', extensions: ['properties'], filenames: ['.env', '.env.local', '.env.development'] },
    { id: 'ignore', name: 'Ignore', extensions: [], filenames: ['.gitignore', '.npmignore', '.dockerignore'] },
    { id: 'lock', name: 'Lock', extensions: ['lock'], filenames: ['package-lock.json', 'yarn.lock', 'Cargo.lock'] }
]

// 获取所有支持的语言（用于下拉菜单）
export const getSupportedLanguages = () =>
    languageConfigs.map(lang => ({
        id: lang.id,
        name: lang.name
    }))

// 根据文件名获取语言
export function getLanguageFromFileName(fileName: string): string {
    const lowerFileName = fileName.toLowerCase()

    // 检查特殊文件名
    for (const lang of languageConfigs) {
        if (lang.filenames?.some(name => name.toLowerCase() === lowerFileName)) {
            return lang.id
        }
    }

    // 检查文件扩展名
    const ext = lowerFileName.split('.').pop() || ''
    const language = languageConfigs.find(lang =>
        lang.extensions.includes(ext)
    )

    return language?.id || 'plaintext'
} 