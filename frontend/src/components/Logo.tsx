interface LogoProps {
  variant?: 'light' | 'dark'
  size?: 'sm' | 'md' | 'lg'
}

const sizes = {
  sm: { icon: 28, text: 'text-xl', sub: 'text-[10px]' },
  md: { icon: 38, text: 'text-3xl', sub: 'text-xs' },
  lg: { icon: 52, text: 'text-4xl', sub: 'text-sm' },
}

export default function Logo({ variant = 'light', size = 'md' }: LogoProps) {
  const s = sizes[size]
  const textColor = variant === 'light' ? 'text-white' : 'text-green-brand'
  const accentColor = variant === 'light' ? 'text-yellow-brand' : 'text-yellow-dark'

  return (
    <div className="flex items-center gap-3">
      <LeafIcon size={s.icon} />
      <div className="flex flex-col leading-none">
        <span className={`${s.text} font-black tracking-tight ${textColor}`}>
          CESI<span className={accentColor}>_</span>
        </span>
        <span className={`${s.sub} font-semibold tracking-[0.2em] uppercase ${accentColor}`}>
          Zen
        </span>
      </div>
    </div>
  )
}

function LeafIcon({ size }: { size: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle cx="20" cy="20" r="20" fill="#f2e2a0" />
      {/* feuille */}
      <path
        d="M20 8 C28 12 32 20 20 32 C8 20 12 12 20 8Z"
        fill="#1a5c32"
      />
      {/* tige */}
      <line x1="20" y1="32" x2="20" y2="36" stroke="#1a5c32" strokeWidth="2" strokeLinecap="round" />
    </svg>
  )
}
