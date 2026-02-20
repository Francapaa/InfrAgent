import Navbar from "./components/landing/Navbar";
import Hero from "./components/landing/Hero";
import Stats from "./components/landing/Stats";
import Features from "./components/landing/Features";
import HowItWorks from "./components/landing/HowItWorks";
import FAQ from "./components/landing/FAQ";
import CTA from "./components/landing/CTA";
import Footer from "./components/landing/Footer";

export default function Home() {
  return (
    <main className="min-h-screen bg-background text-foreground">
      <Navbar />
      <Hero />
      <Stats />
      <Features />
      <HowItWorks />
      <FAQ />
      <CTA />
      <Footer />
    </main>
  );
}
