import type { CapacitorConfig } from '@capacitor/cli'

const config: CapacitorConfig = {
  appId: 'fr.cesizen.app',
  appName: 'CESIZen',
  webDir: 'dist',
  server: {
    // Use https scheme so cookies and localStorage work correctly on Android
    androidScheme: 'https',
  },
  plugins: {
    // Route les requêtes fetch/XHR via la couche HTTP native :
    // évite le blocage "contenu mixte" (app en https://localhost appelant une API en http://)
    // et les soucis de CORS. Indispensable pour joindre le backend de dev en HTTP.
    CapacitorHttp: {
      enabled: true,
    },
  },
}

export default config
