// --- Helper for Icons ---
// Using inline SVGs to keep everything in one file.
const InfoIcon = ({ className = "w-6 h-6" }) => (
  <svg className={className} xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="12" cy="12" r="10" />
    <line x1="12" y1="16" x2="12" y2="12" />
    <line x1="12" y1="8" x2="12.01" y2="8" />
  </svg>
);

const WarningIcon = ({ className = "w-6 h-6" }) => (
  <svg className={className} xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="m21.73 18-8-14a2 2 0 0 0-3.46 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z" />
    <line x1="12" y1="9" x2="12" y2="13" />
    <line x1="12" y1="17" x2="12.01" y2="17" />
  </svg>
);

const TipIcon = ({ className = "w-6 h-6" }) => (
  <svg className={className} xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z" />
    <path d="m9 12 2 2 4-4" />
  </svg>
);

const NoteIcon = ({ className = "w-6 h-6" }) => (
  <svg className={className} xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M4 22h14a2 2 0 0 0 2-2V7.5L14.5 2H6a2 2 0 0 0-2 2v4" />
    <polyline points="14 2 14 8 20 8" />
    <path d="M2 15h10" />
    <path d="M2 19h5" />
  </svg>
);


// --- Admonition Component ---
// This component displays a styled block for notes, warnings, tips, etc.
const Admonition = ({ type = 'note', title, children }) => {
  const styles = {
    note: {
      bgColor: 'bg-blue-50 dark:bg-blue-900/20',
      borderColor: 'border-blue-200 dark:border-blue-500/30',
      iconColor: 'text-blue-500',
      titleColor: 'text-blue-800 dark:text-blue-300',
      icon: <NoteIcon />,
    },
    warning: {
      bgColor: 'bg-red-50 dark:bg-red-900/20',
      borderColor: 'border-red-200 dark:border-red-500/30',
      iconColor: 'text-red-500',
      titleColor: 'text-red-800 dark:text-red-300',
      icon: <WarningIcon />,
    },
    tip: {
      bgColor: 'bg-green-50 dark:bg-green-900/20',
      borderColor: 'border-green-200 dark:border-green-500/30',
      iconColor: 'text-green-600',
      titleColor: 'text-green-800 dark:text-green-300',
      icon: <TipIcon />,
    },
    info: {
      bgColor: 'bg-purple-50 dark:bg-purple-900/20',
      borderColor: 'border-purple-200 dark:border-purple-500/30',
      iconColor: 'text-purple-500',
      titleColor: 'text-purple-800 dark:text-purple-300',
      icon: <InfoIcon />,
    }
  };

  const currentStyle = styles[type] || styles.note;
  const defaultTitle = title || type.charAt(0).toUpperCase() + type.slice(1);

  return (
    <div className={`not-prose mx-auto my-6 flex gap-4 rounded-lg border p-4 max-w-lg ${currentStyle.bgColor} ${currentStyle.borderColor}`}>
      <div className={`mt-1 flex-shrink-0 ${currentStyle.iconColor}`}>
        {currentStyle.icon}
      </div>
      <div className="flex-grow">
        <h3 className={`text-lg font-semibold ${currentStyle.titleColor}`}>{defaultTitle}</h3>
        <div className="prose prose-sm dark:prose-invert max-w-none text-gray-700 dark:text-gray-300 mt-2">
          {children}
        </div>
      </div>
    </div>
  );
};

export default Admonition;