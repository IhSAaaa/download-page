import { render, screen } from '@testing-library/react';
import LoadingSpinner from '../LoadingSpinner';

describe('LoadingSpinner', () => {
  it('renders loading spinner with default text', () => {
    render(<LoadingSpinner />);
    
    const loadingText = screen.getByText('Loading...');
    expect(loadingText).toBeInTheDocument();
  });

  it('renders loading spinner with custom text', () => {
    const customText = 'Please wait...';
    render(<LoadingSpinner text={customText} />);
    
    const loadingText = screen.getByText(customText);
    expect(loadingText).toBeInTheDocument();
  });

  it('applies custom size class', () => {
    render(<LoadingSpinner size="lg" />);
    
    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('w-8', 'h-8');
  });

  it('applies custom color class', () => {
    render(<LoadingSpinner color="blue" />);
    
    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('border-blue-500');
  });
}); 