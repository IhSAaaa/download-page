# URL Shortener Frontend

A modern, elegant, and responsive frontend for the URL Shortener service built with React, TypeScript, and Tailwind CSS.

## ✨ Features

### 🎨 Modern Design
- **Clean & Elegant UI**: Minimalist design with smooth animations and transitions
- **Responsive Layout**: Optimized for desktop, tablet, and mobile devices
- **Glass Morphism**: Modern glass-like effects with backdrop blur
- **Gradient Backgrounds**: Beautiful gradient backgrounds and text effects
- **Smooth Animations**: Fade-in, slide-up, and scale animations for better UX

### 🌙 Dark/Light Mode
- **Theme Toggle**: Elegant toggle switch with smooth transitions
- **System Preference**: Automatically detects user's system theme preference
- **Persistent Storage**: Remembers user's theme choice across sessions
- **Seamless Switching**: Instant theme switching without page reload

### 🚀 Enhanced User Experience
- **Loading States**: Modern loading spinners with contextual messages
- **Hover Effects**: Interactive hover states with scale and shadow effects
- **Toast Notifications**: User-friendly success and error messages
- **Copy to Clipboard**: One-click URL copying with feedback
- **Keyboard Navigation**: Full keyboard accessibility support

### 📊 Dashboard Features
- **Statistics Cards**: Visual overview of total URLs, clicks, and active URLs
- **Modern Cards**: Hover effects and clean typography
- **Pagination**: Smooth pagination with modern controls
- **Action Buttons**: Quick access to copy, analytics, and delete functions

## 🛠️ Technology Stack

- **React 18** - Modern React with hooks and functional components
- **TypeScript** - Type-safe development
- **Tailwind CSS** - Utility-first CSS framework
- **Lucide React** - Beautiful and consistent icons
- **React Router** - Client-side routing
- **React Query** - Data fetching and caching
- **React Hot Toast** - Toast notifications

## 🎯 Design System

### Colors
- **Primary**: Blue gradient (#3B82F6 to #1D4ED8)
- **Secondary**: Gray scale with dark mode variants
- **Accent**: Purple and green for special elements
- **Status**: Red for errors, yellow for warnings, green for success

### Typography
- **Font**: Inter (system fallback)
- **Weights**: Regular (400), Medium (500), Semibold (600), Bold (700)
- **Sizes**: Responsive text scaling

### Components
- **Buttons**: Primary, secondary, and ghost variants with hover effects
- **Cards**: Rounded corners with shadows and hover animations
- **Inputs**: Modern form fields with focus states
- **Navigation**: Clean navigation with active states

## 🚀 Getting Started

### Prerequisites
- Node.js 18+ 
- npm or yarn

### Installation
```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Environment Variables
Create a `.env` file in the frontend directory:
```env
VITE_API_BASE_URL=http://localhost:8080
```

## 🎨 Customization

### Theme Colors
Modify colors in `tailwind.config.js`:
```javascript
colors: {
  primary: {
    50: '#eff6ff',
    // ... other shades
    900: '#1e3a8a',
  },
}
```

### Animations
Add custom animations in `tailwind.config.js`:
```javascript
animation: {
  'fade-in': 'fadeIn 0.5s ease-in-out',
  'slide-up': 'slideInFromBottom 0.5s ease-out',
  // ... more animations
}
```

## 📱 Responsive Design

The application is fully responsive with breakpoints:
- **Mobile**: < 640px
- **Tablet**: 640px - 1024px  
- **Desktop**: > 1024px

## ♿ Accessibility

- **ARIA Labels**: Proper accessibility labels
- **Keyboard Navigation**: Full keyboard support
- **Focus Management**: Clear focus indicators
- **Screen Reader**: Compatible with screen readers
- **Color Contrast**: WCAG AA compliant color ratios

## 🔧 Development

### Project Structure
```
src/
├── components/          # Reusable components
│   ├── Layout.tsx      # Main layout with navigation
│   ├── ThemeToggle.tsx # Dark/light mode toggle
│   └── LoadingSpinner.tsx # Loading component
├── contexts/           # React contexts
│   └── ThemeContext.tsx # Theme management
├── pages/              # Page components
│   ├── Home.tsx        # URL shortening form
│   ├── Dashboard.tsx   # URL management
│   └── Analytics.tsx   # URL analytics
├── services/           # API services
├── types/              # TypeScript types
└── App.tsx             # Main app component
```

### Adding New Components
1. Create component in `src/components/`
2. Use Tailwind classes for styling
3. Add dark mode variants with `dark:` prefix
4. Include proper TypeScript types
5. Add accessibility attributes

### Styling Guidelines
- Use Tailwind utility classes
- Prefer composition over custom CSS
- Include dark mode variants
- Add hover and focus states
- Use consistent spacing and typography

## 🚀 Deployment

### Build for Production
```bash
npm run build
```

### Docker Deployment
The frontend is containerized with Nginx for optimal performance:
```bash
docker build -t url-shortener-frontend .
docker run -p 4000:80 url-shortener-frontend
```

## 🤝 Contributing

1. Follow the existing code style
2. Add TypeScript types for new features
3. Include dark mode support
4. Test on multiple screen sizes
5. Ensure accessibility compliance

## 📄 License

This project is licensed under the MIT License. 