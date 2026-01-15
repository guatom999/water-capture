import { useEffect, useState } from "react"
import { useSearchParams } from 'react-router-dom';
import Footer from "./Footer";
import { getMapMarkerDetailService } from "../services/waterLevelService";
import type { LocationDetail, WaterDetailResponse } from "../types/waterDetail";

const SectionDetail = () => {
    const [searchParams] = useSearchParams();
    const locationId = searchParams.get('location_id');

    const [sectionDetail, setSectionDetail] = useState<WaterDetailResponse | null>(null);

    const getSectionDetail = async () => {
        if (locationId) {
            const response = await getMapMarkerDetailService(Number(locationId));
            setSectionDetail(response as WaterDetailResponse)
        }
    };

    useEffect(() => {
        getSectionDetail();
    }, [locationId]);

    useEffect(() => {
        if (sectionDetail) {
            console.log('Markers data:', sectionDetail.markers);
        }
    }, [sectionDetail]);

    const getDangerColor = (danger: string | null) => {
        switch (danger?.toUpperCase()) {
            case 'CRITICAL':
                return 'bg-red-100 text-red-800 border-red-300';
            case 'DANGER':
                return 'bg-red-100 text-red-800 border-red-300';
            case 'WATCH':
                return 'bg-yellow-100 text-yellow-800 border-yellow-300';
            case 'SAFE':
                return 'bg-green-100 text-green-800 border-green-300';
            default:
                return 'bg-gray-100 text-gray-800 border-gray-300';
        }
    };

    const formatDate = (dateString: string | null) => {
        if (!dateString) return 'N/A';
        return new Date(dateString).toLocaleString('th-TH', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-cyan-50">
            {/* Header */}
            <div className="bg-gradient-to-r from-blue-600 to-cyan-600 text-white py-8 shadow-lg">
                <div className="container mx-auto px-4">
                    <h1 className="text-3xl font-bold mb-2">📊 Water Level Details</h1>
                    <p className="text-blue-100">Location ID: {locationId}</p>
                </div>
            </div>

            {sectionDetail ? (
                <div className="container mx-auto px-4 py-8">
                    {sectionDetail.markers.length === 0 ? (
                        <div className="bg-white rounded-lg shadow-md p-8 text-center">
                            <div className="text-gray-400 text-6xl mb-4">📭</div>
                            <p className="text-gray-600 text-lg">No water level data found for this location.</p>
                        </div>
                    ) : (
                        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                            {sectionDetail.markers.map((marker: LocationDetail) => (
                                <div key={marker.location_id} className="bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-2xl transition-shadow duration-300 border border-gray-100">
                                    {/* Image Section */}
                                    <div className="relative h-48 bg-gradient-to-br from-blue-400 to-cyan-500 overflow-hidden">
                                        <img
                                            src={`http://localhost:8080/images/${marker.image}`}
                                            alt={`Location ${marker.location_id}`}
                                            className="w-full h-full object-cover"
                                            onError={(e) => {
                                                e.currentTarget.src = 'https://via.placeholder.com/400x300?text=No+Image';
                                            }}
                                        />
                                        {/* Danger Badge */}
                                        {marker.danger && (
                                            <div className="absolute top-3 right-3">
                                                <span className={`px-3 py-1 rounded-full text-xs font-semibold border ${getDangerColor(marker.danger)} backdrop-blur-sm`}>
                                                    {marker.danger}
                                                </span>
                                            </div>
                                        )}
                                    </div>

                                    {/* Content Section */}
                                    <div className="p-6">
                                        {/* Water Level - Main Info */}
                                        <div className="mb-4 text-center">
                                            <div className="text-sm text-gray-500 mb-1">Water Level</div>
                                            <div className="text-4xl font-bold text-blue-600">
                                                {marker.level_cm ? `${marker.level_cm.toFixed(2)}` : 'N/A'}
                                            </div>
                                            <div className="text-sm text-gray-500">centimeters</div>
                                        </div>

                                        {/* Flood Status */}
                                        <div className="mb-4 flex items-center justify-center gap-2">
                                            <span className="text-gray-600">Flood Status:</span>
                                            {marker.is_flooded !== null ? (
                                                <span className={`px-3 py-1 rounded-full text-xs font-semibold ${marker.is_flooded ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'}`}>
                                                    {marker.is_flooded ? '🌊 Flooded' : '✅ Normal'}
                                                </span>
                                            ) : (
                                                <span className="text-gray-400">N/A</span>
                                            )}
                                        </div>

                                        {/* Divider */}
                                        <hr className="my-4 border-gray-200" />

                                        {/* Details */}
                                        <div className="space-y-3 text-sm">
                                            <div className="flex items-start gap-2">
                                                <span className="text-gray-500 min-w-[100px]">📅 Measured:</span>
                                                <span className="text-gray-700 font-medium">{formatDate(marker.measured_at)}</span>
                                            </div>

                                            {marker.note && (
                                                <div className="flex items-start gap-2">
                                                    <span className="text-gray-500 min-w-[100px]">📝 Note:</span>
                                                    <span className="text-gray-700">{marker.note}</span>
                                                </div>
                                            )}

                                            {marker.source?.Valid && (
                                                <div className="flex items-start gap-2">
                                                    <span className="text-gray-500 min-w-[100px]">🔗 Source:</span>
                                                    <span className="text-gray-700">{marker.source.String}</span>
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            ) : (
                <div className="flex items-center justify-center min-h-[60vh]">
                    <div className="text-center">
                        <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-blue-600 mx-auto mb-4"></div>
                        <p className="text-gray-600 text-lg">Loading water level data...</p>
                    </div>
                </div>
            )}

            <Footer />
        </div>
    );
};

export default SectionDetail;