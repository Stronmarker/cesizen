import { Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import AppShell from '../components/AppShell'
import Logo from '../components/Logo'

function greeting() {
  const h = new Date().getHours()
  if (h < 12) return 'Bonjour'
  if (h < 18) return 'Bon après-midi'
  return 'Bonsoir'
}

export default function HomePage() {
  const { user } = useAuth()
  const name = user?.first_name || user?.nickname || user?.email?.split('@')[0] || 'vous'

  return (
    <AppShell>
      {/* Header */}
      <div className="relative px-6 pt-[calc(2.5rem_+_env(safe-area-inset-top))] pb-8 overflow-hidden">
        <div className="absolute -top-12 -right-12 w-52 h-52 rounded-full bg-green-brand/8 pointer-events-none" />
        <div className="absolute top-16 right-8 w-20 h-20 rounded-full bg-yellow-brand/15 pointer-events-none" />

        <Logo size="sm" variant="dark" />

        <div className="mt-7">
          <p className="text-green-brand/45 text-sm font-medium">{greeting()} 👋</p>
          <h1 className="text-green-brand text-3xl font-black mt-0.5 leading-tight">{name}</h1>
          <p className="text-green-brand/45 text-sm mt-1.5">Comment vous sentez-vous aujourd'hui ?</p>
        </div>
      </div>

      <div className="px-4 pb-4 flex flex-col gap-3">
        <QuickAction
          title="Tracker mes émotions"
          description="Enregistrez votre ressenti du moment"
          icon="💚"
          href="/tracker"
        />
        <QuickAction
          title="Articles & conseils"
          description="Ressources sur la santé mentale"
          icon="📖"
          href="/info"
        />
        <QuickAction
          title="Mon profil"
          description="Gérer mon compte"
          icon="👤"
          href="/profile"
        />
      </div>
    </AppShell>
  )
}

function QuickAction({
  title, description, icon, href,
}: Readonly<{ title: string; description: string; icon: string; href: string }>) {
  return (
    <Link
      to={href}
      className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-4 flex items-center gap-4
                 active:scale-[0.98] active:opacity-90 transition-all shadow-sm"
    >
      <div className="w-12 h-12 bg-green-brand/12 rounded-xl flex items-center justify-center text-2xl flex-shrink-0">
        {icon}
      </div>
      <div className="flex-1 min-w-0">
        <p className="font-semibold text-green-brand text-sm">{title}</p>
        <p className="text-green-brand/50 text-xs mt-0.5">{description}</p>
      </div>
      <svg className="w-4 h-4 text-green-brand/30 shrink-0" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
      </svg>
    </Link>
  )
}
