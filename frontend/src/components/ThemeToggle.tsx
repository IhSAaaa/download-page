import { Moon, Sun } from 'lucide-react'
import { useTheme } from '@/contexts/ThemeContext'

const ThemeToggle = () => {
  const { theme, toggleTheme } = useTheme()

  return (
    <button
      onClick={toggleTheme}
      className="theme-toggle h-8 w-14 bg-gradient-to-r from-gray-100 to-gray-200 dark:from-gray-700 dark:to-gray-800 border border-gray-200 dark:border-gray-600 shadow-sm hover:shadow-md"
      aria-label="Toggle theme"
    >
      <div className="relative flex h-6 w-12 items-center justify-between rounded-full bg-white dark:bg-gray-800 px-1">
        <div
          className={`theme-toggle-thumb top-0.5 h-5 w-5 bg-gradient-to-r from-yellow-400 to-orange-500 ${
            theme === 'dark' ? 'translate-x-6' : 'translate-x-0'
          }`}
        />
        <Sun className={`theme-toggle-icon h-3 w-3 ${
          theme === 'light' ? 'text-yellow-600' : 'text-gray-400'
        }`} />
        <Moon className={`theme-toggle-icon h-3 w-3 ${
          theme === 'dark' ? 'text-blue-400' : 'text-gray-400'
        }`} />
      </div>
    </button>
  )
}

export default ThemeToggle 