import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from 'react-query'
import { toast } from 'react-hot-toast'
import { Link as LinkIcon, Copy, Trash2, BarChart3, ExternalLink, Sparkles, Calendar, Tag, FileText } from 'lucide-react'
import { Link } from 'react-router-dom'
import { urlService } from '@/services/api'
import { CreateURLRequest, URLResponse } from '@/types'

const Home = () => {
  const queryClient = useQueryClient()
  const [form, setForm] = useState<CreateURLRequest>({
    original_url: '',
    custom_code: '',
    title: '',
    description: '',
    expires_at: '',
  })

  // Query for recent URLs
  const { data: urlsData, isLoading: loadingUrls } = useQuery(
    ['urls', 1, 10],
    () => urlService.getAllURLs(1, 10)
  )

  // Mutation for creating short URL
  const createURLMutation = useMutation(urlService.createShortURL, {
    onSuccess: () => {
      toast.success('URL shortened successfully!')
      setForm({
        original_url: '',
        custom_code: '',
        title: '',
        description: '',
        expires_at: '',
      })
      queryClient.invalidateQueries(['urls'])
    },
    onError: (error: any) => {
      const errorMessage = error.response?.data?.error || error.response?.data?.details || 'Failed to shorten URL'
      toast.error(errorMessage)
    },
  })

  // Mutation for deleting URL
  const deleteURLMutation = useMutation(urlService.deleteURL, {
    onSuccess: () => {
      toast.success('URL deleted successfully!')
      queryClient.invalidateQueries(['urls'])
    },
    onError: () => {
      toast.error('Failed to delete URL')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.original_url) return

    const payload: CreateURLRequest = {
      original_url: form.original_url,
      custom_code: form.custom_code || undefined,
      title: form.title || undefined,
      description: form.description || undefined,
    }
    
    // Only include expires_at if it's not empty
    if (form.expires_at && form.expires_at.trim() !== '') {
      payload.expires_at = new Date(form.expires_at).toISOString()
    }

    createURLMutation.mutate(payload)
  }

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
      toast.success('Copied to clipboard!')
    } catch (error) {
      toast.error('Failed to copy to clipboard')
    }
  }

  const handleDelete = (id: string) => {
    if (window.confirm('Are you sure you want to delete this URL?')) {
      deleteURLMutation.mutate(id)
    }
  }

  return (
    <div className="space-y-12">
      {/* Hero Section */}
      <div className="text-center slide-up">
        <div className="inline-flex items-center px-4 py-2 rounded-full bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 text-sm font-medium mb-6">
          <Sparkles className="h-4 w-4 mr-2" />
          Modern URL Shortener
        </div>
        <h1 className="text-5xl md:text-6xl font-bold bg-gradient-to-r from-gray-900 via-blue-800 to-purple-800 dark:from-gray-100 dark:via-blue-400 dark:to-purple-400 bg-clip-text text-transparent mb-6">
          Shorten Your URLs
        </h1>
        <p className="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto leading-relaxed">
          Create short, memorable links and track their performance with detailed analytics. 
          <span className="block mt-2 text-lg text-gray-500 dark:text-gray-400">
            Fast, secure, and beautiful.
          </span>
        </p>
      </div>

      {/* URL Shortener Form */}
      <div className="max-w-3xl mx-auto scale-in">
        <div className="card card-hover">
          <div className="flex items-center mb-6">
            <div className="h-10 w-10 bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl flex items-center justify-center mr-4">
              <LinkIcon className="h-5 w-5 text-white" />
            </div>
            <div>
              <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">Create Short URL</h2>
              <p className="text-gray-500 dark:text-gray-400">Fill in the details below</p>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="originalUrl" className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">
                Original URL *
              </label>
              <div className="relative">
                <input
                  type="url"
                  id="originalUrl"
                  value={form.original_url}
                  onChange={(e) => setForm({ ...form, original_url: e.target.value })}
                  placeholder="https://example.com/very-long-url"
                  className="input-field pl-12"
                  required
                />
                                 <LinkIcon className="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label htmlFor="customCode" className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">
                  Custom Code (Optional)
                </label>
                <div className="relative">
                  <input
                    type="text"
                    id="customCode"
                    value={form.custom_code || ''}
                    onChange={(e) => setForm({ ...form, custom_code: e.target.value })}
                    placeholder="my-custom-link"
                    className="input-field pl-12"
                  />
                  <Tag className="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
                </div>
              </div>
              <div>
                <label htmlFor="expiresAt" className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">
                  Expires At (Optional)
                </label>
                <div className="relative">
                  <input
                    type="datetime-local"
                    id="expiresAt"
                    value={form.expires_at || ''}
                    onChange={(e) => setForm({ ...form, expires_at: e.target.value })}
                    className="input-field pl-12"
                  />
                  <Calendar className="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
                </div>
              </div>
            </div>

            <div>
              <label htmlFor="title" className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">
                Title (Optional)
              </label>
              <input
                type="text"
                id="title"
                value={form.title || ''}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
                placeholder="My Awesome Link"
                className="input-field"
              />
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">
                Description (Optional)
              </label>
              <div className="relative">
                <textarea
                  id="description"
                  value={form.description || ''}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  placeholder="Brief description of this link"
                  rows={3}
                  className="input-field pl-12 resize-none"
                />
                <FileText className="absolute left-4 top-4 h-5 w-5 text-gray-400" />
              </div>
            </div>

            <button
              type="submit"
              disabled={createURLMutation.isLoading}
              className="btn-primary w-full flex items-center justify-center space-x-2"
            >
              {createURLMutation.isLoading ? (
                <>
                  <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                  <span>Creating...</span>
                </>
              ) : (
                <>
                                     <LinkIcon className="h-5 w-5" />
                   <span>Shorten URL</span>
                </>
              )}
            </button>
          </form>
        </div>
      </div>

      {/* Recent URLs Section */}
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Recent URLs</h2>
          <div className="h-px flex-1 bg-gradient-to-r from-transparent via-gray-300 dark:via-gray-600 to-transparent mx-4"></div>
        </div>

        {loadingUrls ? (
          <div className="card text-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p className="text-gray-500 dark:text-gray-400">Loading your URLs...</p>
          </div>
        ) : urlsData?.data && urlsData.data.length > 0 ? (
          <div className="grid gap-4">
            {urlsData.data.map((url: URLResponse) => (
              <div key={url.id} className="card card-hover group">
                <div className="flex items-center justify-between">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center space-x-3 mb-2">
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 truncate">
                        {url.title || 'Untitled Link'}
                      </h3>
                      {url.expires_at && (
                        <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-400">
                          Expires {new Date(url.expires_at).toLocaleDateString()}
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mb-2 truncate">
                      {url.original_url}
                    </p>
                    <div className="flex items-center space-x-4">
                      <a
                        href={url.short_url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-blue-600 dark:text-blue-400 font-medium hover:underline flex items-center space-x-1"
                      >
                        <span>{url.short_url}</span>
                        <ExternalLink className="h-3 w-3" />
                      </a>
                      <span className="text-gray-400">•</span>
                      <span className="text-sm text-gray-500 dark:text-gray-400">
                        Created {new Date(url.created_at).toLocaleDateString()}
                      </span>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button
                      onClick={() => copyToClipboard(url.short_url)}
                      className="btn-ghost p-2"
                      title="Copy URL"
                    >
                      <Copy className="h-4 w-4" />
                    </button>
                    <Link
                      to={`/analytics/${url.id}`}
                      className="btn-ghost p-2"
                      title="View Analytics"
                    >
                      <BarChart3 className="h-4 w-4" />
                    </Link>
                    <button
                      onClick={() => handleDelete(url.id)}
                      className="btn-ghost p-2 text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
                      title="Delete URL"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="card text-center py-16">
                         <div className="w-24 h-24 bg-gray-100 dark:bg-gray-800 rounded-full flex items-center justify-center mx-auto mb-6">
               <LinkIcon className="h-12 w-12 text-gray-400" />
             </div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2">
              No URLs created yet
            </h3>
            <p className="text-gray-500 dark:text-gray-400">
              Create your first short URL above!
            </p>
          </div>
        )}
      </div>
    </div>
  )
}

export default Home 