import Footer from '../components/Footer';
import Header from '../components/Header';
import MapSection from '../components/MapSection';
// import MapSectionOld from "../components/MapSectionOld"

const LandingPage = () => {
    return (
        <div className="min-h-screen bg-white">
            <MapSection />
            <Footer />
        </div>
    );
};

export default LandingPage;
