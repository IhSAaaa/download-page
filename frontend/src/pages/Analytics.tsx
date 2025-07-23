import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { urlService } from '@/services/api'
import { Analytics as AnalyticsType } from '@/types'

const Analytics = () => {
  const { id } = useParams<{ id: string }>()
  const [analytics, setAnalytics] = useState<AnalyticsType | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (id) {
      loadAnalytics()
    }
  }, [id])

  const loadAnalytics = async () => {
    try {
      setLoading(true)
      const data = await urlService.getURLAnalytics(id!)
      setAnalytics(data)
    } catch (err) {
      setError('Failed to load analytics')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (error || !analytics) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error || 'Analytics not found'}
        </div>
        <Link
          to="/dashboard"
          className="mt-4 inline-block bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          Back to Dashboard
        </Link>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <Link
          to="/dashboard"
          className="text-blue-600 hover:text-blue-800 font-medium"
        >
          ← Back to Dashboard
        </Link>
        <h1 className="text-3xl font-bold text-gray-900 mt-4">URL Analytics</h1>
      </div>

      {/* Overview Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-medium text-gray-900">Total Clicks</h3>
          <p className="text-3xl font-bold text-blue-600">{analytics.total_clicks}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-medium text-gray-900">Unique Clicks</h3>
          <p className="text-3xl font-bold text-green-600">{analytics.unique_clicks}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-medium text-gray-900">Last Click</h3>
          <p className="text-lg text-gray-600">
            {analytics.last_clicked_at 
              ? new Date(analytics.last_clicked_at).toLocaleDateString()
              : 'Never'
            }
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Top Countries */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Top Countries</h3>
          {analytics.top_countries && analytics.top_countries.length > 0 ? (
            <div className="space-y-3">
              {analytics.top_countries.map((country, index) => (
                <div key={index} className="flex justify-between items-center">
                  <span className="text-gray-700">{country.country}</span>
                  <span className="font-medium text-blue-600">{country.clicks}</span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500">No country data available</p>
          )}
        </div>

        {/* Top Devices */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Top Devices</h3>
          {analytics.top_devices && analytics.top_devices.length > 0 ? (
            <div className="space-y-3">
              {analytics.top_devices.map((device, index) => (
                <div key={index} className="flex justify-between items-center">
                  <span className="text-gray-700">{device.device}</span>
                  <span className="font-medium text-green-600">{device.clicks}</span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500">No device data available</p>
          )}
        </div>

        {/* Top Browsers */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Top Browsers</h3>
          {analytics.top_browsers && analytics.top_browsers.length > 0 ? (
            <div className="space-y-3">
              {analytics.top_browsers.map((browser, index) => (
                <div key={index} className="flex justify-between items-center">
                  <span className="text-gray-700">{browser.browser}</span>
                  <span className="font-medium text-purple-600">{browser.clicks}</span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500">No browser data available</p>
          )}
        </div>

        {/* Click Timeline */}
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Click Timeline (Last 30 Days)</h3>
          {analytics.click_timeline && analytics.click_timeline.length > 0 ? (
            <div className="space-y-3">
              {analytics.click_timeline.map((timeline, index) => (
                <div key={index} className="flex justify-between items-center">
                  <span className="text-gray-700">
                    {new Date(timeline.date).toLocaleDateString()}
                  </span>
                  <span className="font-medium text-orange-600">{timeline.clicks}</span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500">No timeline data available</p>
          )}
        </div>
      </div>
    </div>
  )
}

export default Analytics 