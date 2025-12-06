import React from 'react';

const center = {
    lat: 13.7563, // Bangkok coordinates
    lng: 100.5018
}

const MapSection = () => {
    return (
        <section className="py-16 bg-gray-50">
            <div className="container mx-auto px-4">
                <div className="text-center mb-12">
                    <h2 className="text-3xl font-bold text-gray-900 mb-4">Find Us</h2>
                    <p className="text-gray-600 max-w-2xl mx-auto">
                        Visit our office in the heart of the city. We are easily accessible by public transport.
                    </p>
                </div>

                <div className="rounded-xl overflow-hidden shadow-lg border border-gray-200 h-[400px]">
                    <iframe
                        src={`https://maps.google.com/maps?q=${center.lat},${center.lng}&z=15&output=embed`}
                        width="100%"
                        height="100%"
                        style={{ border: 0 }}
                        allowFullScreen={true}
                        loading="lazy"
                        referrerPolicy="no-referrer-when-downgrade"
                        title="Google Map"
                    ></iframe>
                </div>
            </div>
        </section >
    );
};

export default MapSection;
