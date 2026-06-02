import { cn } from '../../lib/utils';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary';
  isLoading?: boolean;
}

export function Button({ variant = 'primary', isLoading, className, children, disabled, ...props }: ButtonProps) {
  return (
    <button
      disabled={disabled || isLoading}
      className={cn(
        'inline-flex items-center justify-center gap-2 rounded-lg px-4 py-2.5',
        'text-sm font-medium transition-all duration-150',
        'focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
        'disabled:pointer-events-none disabled:opacity-50',
        variant === 'primary' && [
          'bg-forest text-white',
          'hover:bg-forest-dark',
          'focus-visible:ring-forest',
        ],
        variant === 'secondary' && [
          'bg-surface text-ink border border-border',
          'hover:bg-surface-hover',
          'focus-visible:ring-border-strong',
        ],
        className,
      )}
      {...props}
    >
      {isLoading && (
        <span className="h-3.5 w-3.5 rounded-full border-2 border-current border-t-transparent animate-spin" />
      )}
      {children}
    </button>
  );
}
