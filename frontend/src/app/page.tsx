import Header from '@/components/Header';

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col bg-[#050505] text-[#f0f0f0]">
      <Header />
      <main className="flex-1 max-w-5xl mx-auto px-6 py-16 flex flex-col items-center justify-center text-center">
        <h1 className="text-4xl font-extrabold tracking-tight mb-4 text-[#ff1e1e]">
          YouTube Data Tool
        </h1>
        <p className="text-[#888888] text-lg max-w-2xl mb-8">
          A simple tool to pull video details, descriptions, subtitles, and transcripts from YouTube links or playlists.
        </p>
        <div className="flex gap-4">
          <a
            href="/export"
            className="bg-[#ff1e1e] text-white font-semibold px-6 py-3 rounded hover:bg-[#e01212] transition-colors"
          >
            Go to Export Tool
          </a>
        </div>
      </main>
    </div>
  );
}
