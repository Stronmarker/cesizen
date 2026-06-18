import { NavLink } from 'react-router-dom'

interface AppShellProps {
  children: React.ReactNode
}

export default function AppShell({ children }: Readonly<AppShellProps>) {
  return (
    <div className="min-h-screen flex flex-col bg-transparent">
      <main className="flex-1 overflow-y-auto pb-24">
        {children}
      </main>
      <nav
        className="fixed bottom-0 inset-x-0 bg-cream/70 backdrop-blur-xl border-t border-white/50 shadow-[0_-8px_24px_rgba(26,92,50,0.06)] flex"
        style={{ paddingBottom: 'env(safe-area-inset-bottom)' }}
      >
        <NavItem to="/" icon={<HomeIcon />} label="Accueil" />
        <NavItem to="/tracker" icon={<HeartIcon />} label="Émotions" />
        <NavItem to="/info" icon={<InfoIcon />} label="Infos" />
        <NavItem to="/profile" icon={<UserIcon />} label="Profil" />
      </nav>
    </div>
  )
}

function NavItem({ to, icon, label }: Readonly<{ to: string; icon: React.ReactNode; label: string }>) {
  return (
    <NavLink
      to={to}
      end
      className={({ isActive }) =>
        `flex-1 flex flex-col items-center justify-center py-2.5 gap-0.5 text-xs font-medium transition-colors ${
          isActive ? 'text-green-brand' : 'text-green-brand/35'
        }`
      }
    >
      {icon}
      {label}
    </NavLink>
  )
}

function HomeIcon() {
  return (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
    </svg>
  )
}

function HeartIcon() {
  return (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
    </svg>
  )
}

function InfoIcon() {
  return (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  )
}

function UserIcon() {
  return (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
    </svg>
  )
}
