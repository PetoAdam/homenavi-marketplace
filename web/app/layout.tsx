import './globals.css';

export const metadata = {
  title: 'Homenavi Marketplace',
  description: 'Browse and manage Homenavi integrations.',
  icons: {
    icon: '/icon.png'
  }
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="marketplace-bg min-h-screen">
        {children}
      </body>
    </html>
  );
}
