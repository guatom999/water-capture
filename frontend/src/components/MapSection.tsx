// import React, { useState } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L, { Handler } from 'leaflet';
import icon from 'leaflet/dist/images/marker-icon.png';
import iconShadow from 'leaflet/dist/images/marker-shadow.png';
import Header from './Header';

// let DefaultIcon = L.icon({
//     iconUrl: icon,
//     shadowUrl: iconShadow,
//     iconSize: [25, 41],
//     iconAnchor: [12, 41]
// });

// L.Marker.prototype.options.icon = DefaultIcon;

const locations = [
    { id: 1, name: "สำนักงานใหญ่ (Bangkok)", lat: 13.7563, lng: 100.5018 },
    { id: 2, name: "สาขาเชียงใหม่", lat: 18.7883, lng: 98.9853 },
    { id: 3, name: "สาขาภูเก็ต", lat: 7.8804, lng: 98.3923 },
    { id: 4, name: "สาขารังสิต", lat: 13.98691942837372, lng: 100.61850610106555 },
];


const MapSection = () => {

    return (
        <section className="bg-gray-50 w-full h-screen relative">
            <div
                className='w-full h-full'
            >
                <MapContainer
                    center={[13.7563, 100.5018]}
                    zoom={10}
                    scrollWheelZoom={true}
                    style={{ height: '100%', width: '100%' }}
                >
                    {/* TileLayer คือภาพแผนที่ (ใช้ของ OpenStreetMap ฟรี) */}
                    <TileLayer
                        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    />

                    {/* วนลูปสร้าง Marker ตามข้อมูลใน locations */}
                    {locations.map(loc => (
                        <Marker key={loc.id} position={[loc.lat, loc.lng]}>
                            <Popup>
                                <div className="font-bold">{loc.name}</div>
                                <p>รายละเอียดสาขา...</p>
                                <img src="https://scontent.fbkk29-7.fna.fbcdn.net/v/t39.30808-6/560614609_1217327910441902_6428355889022782810_n.jpg?_nc_cat=108&ccb=1-7&_nc_sid=833d8c&_nc_ohc=lwu-Echj39kQ7kNvwHPDyCy&_nc_oc=AdnznDahG3A0c3_xwWmSSkJ5TXYhElY0U9F4dSt7pN9yieIaz4OwrijIagdwX3WLxNA&_nc_zt=23&_nc_ht=scontent.fbkk29-7.fna&_nc_gid=Jx2f49ujdcB0fG5hy4UhWw&oh=00_AflfkHwb1oTehjOFJ-QuQeSIu9zZj-KEjVe-V81BcjKLgQ&oe=6939FB5C"></img>
                            </Popup>
                        </Marker>
                    ))}
                </MapContainer>
            </div>
        </section>
    );
};

export default MapSection;