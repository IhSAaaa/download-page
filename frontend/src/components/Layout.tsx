import { ReactNode } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { Link as LinkIcon, BarChart3, Home } from 'lucide-react'
import ThemeToggle from './ThemeToggle'

interface LayoutProps {
  children: ReactNode
}

const Layout = ({ children }: LayoutProps) => {
  const location = useLocation()

  const navigation = [
    { name: 'Home', href: '/', icon: Home },
    { name: 'Dashboard', href: '/dashboard', icon: BarChart3 },
  ]

  return (
    <div className="min-h-screen gradient-bg transition-colors duration-300">
      {/* Header */}
      <header className="glass sticky top-0 z-50 backdrop-blur-lg border-b border-gray-200/50 dark:border-gray-700/50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <div className="relative">
                <LinkIcon className="h-8 w-8 text-blue-600 dark:text-blue-400 transition-colors duration-300" />
                <div className="absolute -top-1 -right-1 h-3 w-3 bg-green-500 rounded-full animate-pulse"></div>
              </div>
              <h1 className="ml-3 text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-600 dark:from-gray-100 dark:to-gray-400 bg-clip-text text-transparent">
                URL Shortener
              </h1>
            </div>
            
            <div className="flex items-center space-x-4">
              <nav className="flex space-x-2">
                {navigation.map((item) => {
                  const Icon = item.icon
                  const isActive = location.pathname === item.href
                  
                  return (
                    <Link
                      key={item.name}
                      to={item.href}
                      className={`nav-link ${
                        isActive ? 'nav-link-active' : 'nav-link-inactive'
                      }`}
                    >
                      <Icon className="h-4 w-4 mr-2" />
                      {item.name}
                    </Link>
                  )
                })}
              </nav>
              
              <div className="h-6 w-px bg-gray-300 dark:bg-gray-600"></div>
              
              <ThemeToggle />
            </div>
          </div>
        </div>
      </header>

      {/* Main content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="fade-in">
          {children}
        </div>
      </main>

      {/* Footer */}
      <footer className="glass border-t border-gray-200/50 dark:border-gray-700/50 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="text-center">
            <p className="text-gray-600 dark:text-gray-400 text-sm">
              Built with ❤️ using <span className="font-semibold text-blue-600 dark:text-blue-400">Go</span> and <span className="font-semibold text-blue-600 dark:text-blue-400">React</span>
            </p>
          </div>
        </div>
      </footer>
    </div>
  )
}

export default Layout 