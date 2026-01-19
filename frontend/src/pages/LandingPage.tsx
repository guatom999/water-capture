import Footer from '../components/Footer';
import Header from '../components/Header';
import MapSection from '../components/MapSection';

const LandingPage = () => {
    return (
        <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
            <Header />
            <MapSection />
            <Footer />
        </div>
    );
};

export default LandingPage;
