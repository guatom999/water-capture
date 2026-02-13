// Note (ไทย): secret ที่มาจาก config ใช้เฉพาะบนฝั่งเซิร์ฟเวอร์เพื่อเซ็นและตรวจสอบ access token / refresh token.
// ไม่ควรนำ secret นี้ไปไว้หรือเรียกใช้ในโค้ดฝั่งไคลเอนต์ (frontend) เพราะจะเป็นความเสี่ยงด้านความปลอดภัย.
// Frontend ควรเก็บแค่ token ที่ได้รับจากเซิร์ฟเวอร์ (เช่น accessToken/refreshToken) และส่งไปยัง API เมื่อจำเป็น.

// import React from 'react';
// import { Link } from 'react-router-dom';
import Header from '../components/Header';

const user = {
    id: 1,
    email: 'jest_trust@hotmail.com',
    name: 'test',
    role: 'USER',
    subscribed_channel: 'line',
    subscribe_province_id: 13,
    subscribe_location_id: 37,
};

export default function UserProfilePage() {
    return (
        <>
            <Header />
            <div className="min-h-screen bg-gradient-to-br from-gray-100 via-blue-50 to-gray-100">
                <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 pt-10">
                    <div className="bg-white rounded-2xl shadow-xlmb-6">
                        <div className="p-6 sm:p-8">
                            <div className="flex flex-col sm:flex-row items-center sm:items-end gap-5">
                                <div className="w-28 h-28 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center shadow-lg ring-4 ring-white -mt-20 sm:-mt-16">
                                    <span className="text-4xl font-bold text-white uppercase">
                                        {user.name.charAt(0)}
                                    </span>
                                </div>
                                <div className="flex-1 text-center sm:text-left">
                                    <div className="flex flex-col sm:flex-row items-center gap-2">
                                        <h1 className="text-2xl font-bold text-gray-900">{user.name}</h1>
                                        <span className="px-3 py-0.5 text-xs font-semibold rounded-full bg-blue-100 text-blue-700 uppercase tracking-wide">
                                            {user.role}
                                        </span>
                                    </div>
                                    <p className="text-gray-500 mt-1">{user.email}</p>
                                </div>
                                {/* <Link
                                    to="/user-profile/edit"
                                    className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200 shadow hover:shadow-md"
                                >
                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                    </svg>
                                    แก้ไขโปรไฟล์
                                </Link> */}
                            </div>
                        </div>
                    </div>

                    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 py-10">
                        {/* Main Content */}
                        <main className="lg:col-span-2 space-y-6">
                            {/* Personal Info */}
                            <div className="bg-white rounded-2xl shadow-md p-6">
                                <h2 className="text-lg font-bold text-gray-900 flex items-center gap-2 mb-5">
                                    <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                                    </svg>
                                    ข้อมูลส่วนตัว
                                </h2>
                                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                                    {[
                                        { label: 'ID', value: user.id },
                                        { label: 'ชื่อ', value: user.name },
                                        { label: 'อีเมล', value: user.email },
                                        { label: 'บทบาท', value: user.role },
                                    ].map((item) => (
                                        <div key={item.label} className="bg-gray-50 rounded-xl p-4">
                                            <p className="text-xs text-gray-400 uppercase tracking-wider font-semibold">{item.label}</p>
                                            <p className="text-gray-800 font-medium mt-1">{item.value}</p>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            {/* Subscription Info */}
                            <div className="bg-white rounded-2xl shadow-md p-6">
                                <div className="flex items-center gap-2 justify-between gap-2 mb-5 ">
                                    <h2 className="text-lg font-bold text-gray-900 flex items-center ">
                                        <svg className="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                                        </svg>
                                        การแจ้งเตือน
                                    </h2>
                                    <div>
                                        Subscribed 1/1
                                    </div>
                                </div>

                                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                                    <div className="bg-green-50 border border-green-100 rounded-xl p-4 text-center">
                                        <p className="text-xs text-green-500 uppercase tracking-wider font-semibold">ช่องทาง</p>
                                        <p className="text-green-700 font-bold text-lg mt-1 capitalize">{user.subscribed_channel}</p>
                                    </div>
                                    <div className="bg-blue-50 border border-blue-100 rounded-xl p-4 text-center">
                                        <p className="text-xs text-blue-500 uppercase tracking-wider font-semibold">จังหวัด ID</p>
                                        <p className="text-blue-700 font-bold text-lg mt-1">{user.subscribe_province_id}</p>
                                    </div>
                                    <div className="bg-indigo-50 border border-indigo-100 rounded-xl p-4 text-center">
                                        <p className="text-xs text-indigo-500 uppercase tracking-wider font-semibold">สถานี ID</p>
                                        <p className="text-indigo-700 font-bold text-lg mt-1">{user.subscribe_location_id}</p>
                                    </div>
                                </div>
                            </div>
                        </main>

                        {/* Sidebar */}
                        {/* <aside className="space-y-6">
                            <div className="bg-white rounded-2xl shadow-md p-5">
                                <h3 className="text-sm font-bold text-gray-900 mb-3">เมนูลัด</h3>
                                <div className="flex flex-col gap-2">
                                    {[
                                        { to: '/', label: 'หน้าหลัก', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-4 0a2 2 0 110-4m0 4a2 2 0 100-4m0 4v1m0-1a2 2 0 01-2-2V8' },
                                        { to: '/map', label: 'แผนที่', icon: 'M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7' },
                                        { to: '/reports', label: 'รายงาน', icon: 'M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z' },
                                    ].map((item) => (
                                        <Link
                                            key={item.to}
                                            to={item.to}
                                            className="flex items-center gap-3 text-sm text-gray-700 hover:text-blue-700 hover:bg-blue-50 px-3 py-2.5 rounded-xl transition-all duration-200"
                                        >
                                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={item.icon} />
                                            </svg>
                                            {item.label}
                                        </Link>
                                    ))}
                                </div>
                            </div>

                            <div className="bg-gradient-to-br from-blue-600 to-indigo-700 rounded-2xl shadow-md p-5 text-white">
                                <h3 className="text-sm font-bold mb-2">💧 Water Level Monitor</h3>
                                <p className="text-xs text-blue-100 leading-relaxed">
                                    ระบบจะแจ้งเตือนผ่าน <span className="font-semibold capitalize">{user.subscribed_channel}</span> เมื่อระดับน้ำในสถานีของคุณเปลี่ยนแปลง
                                </p>
                            </div>
                        </aside> */}
                    </div>
                </div>
            </div>
        </>
    );
}
