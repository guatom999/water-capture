import Header from '../components/Header';
import MapSection from '../components/MapSection';
// import MapSectionOld from "../components/MapSectionOld"

const LandingPage = () => {
    return (
        <div className="min-h-screen bg-white">
            {/* Map Section */}
            <MapSection />
            {/* <MapSectionOld /> */}

            {/* Footer */}
            <footer className="bg-gray-900 text-white py-8">
                <div className="container mx-auto px-4 text-center">
                    <p>&copy; {new Date().getFullYear()} Your Company. All rights reserved.</p>
                </div>
            </footer>
        </div>
    );
};

export default LandingPage;
