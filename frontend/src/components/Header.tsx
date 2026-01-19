import { useState } from 'react';
import { Link } from 'react-router-dom';

const Header = () => {
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    return (
        <header className="bg-gradient-to-r from-blue-600 via-blue-700 to-indigo-800 shadow-lg sticky top-0 z-50">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    <div className="flex items-center space-x-3">
                        <div className="w-10 h-10 bg-white/20 rounded-lg flex items-center justify-center backdrop-blur-sm">
                            <svg
                                className="w-6 h-6 text-white"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                                />
                            </svg>
                        </div>
                        <div>
                            <h1 className="text-xl font-bold text-white tracking-tight">
                                Water Level Monitor
                            </h1>
                            <p className="text-xs text-blue-200 hidden sm:block">
                                ระบบติดตามระดับน้ำ
                            </p>
                        </div>
                    </div>

                    <nav className="hidden md:flex items-center space-x-6">
                        <Link
                            to="/#"
                            className="text-white/90 hover:text-white transition-colors duration-200 text-sm font-medium"
                        >
                            หน้าหลัก
                        </Link>
                        <Link
                            to="/#"
                            className="text-white/90 hover:text-white transition-colors duration-200 text-sm font-medium"
                        >
                            แผนที่
                        </Link>
                        <a
                            href="#"
                            className="text-white/90 hover:text-white transition-colors duration-200 text-sm font-medium"
                        >
                            รายงาน
                        </a>
                    </nav>

                    <div className="flex items-center space-x-4">
                        <Link
                            to="/login"
                            className="bg-white text-blue-700 px-5 py-2 rounded-lg font-semibold text-sm 
                                       hover:bg-blue-50 transition-all duration-200 shadow-md hover:shadow-lg
                                       transform hover:-translate-y-0.5 active:translate-y-0"
                        >
                            <span className="flex items-center space-x-2">
                                <svg
                                    className="w-4 h-4"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"
                                    />
                                </svg>
                                <span>เข้าสู่ระบบ</span>
                            </span>
                        </Link>

                        <button
                            className="md:hidden text-white p-2 rounded-lg hover:bg-white/10 transition-colors"
                            onClick={() => setIsMenuOpen(!isMenuOpen)}
                        >
                            <svg
                                className="w-6 h-6"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                {isMenuOpen ? (
                                    <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M6 18L18 6M6 6l12 12"
                                    />
                                ) : (
                                    <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M4 6h16M4 12h16M4 18h16"
                                    />
                                )}
                            </svg>
                        </button>
                    </div>
                </div>

                {isMenuOpen && (
                    <div className="md:hidden py-4 border-t border-white/10">
                        <nav className="flex flex-col space-y-3">
                            <a
                                href="/"
                                className="text-white/90 hover:text-white transition-colors duration-200 text-sm font-medium px-2 py-2 rounded-lg hover:bg-white/10"
                            >
                                หน้าหลัก
                            </a>
                            <a
                                href="#"
                                className="text-white/90 hover:text-white transition-colors duration-200 text-sm font-medium px-2 py-2 rounded-lg hover:bg-white/10"
                            >
                                แผนที่
                            </a>
                            <a
                                href="#"
                                className="text-white/90 hover:text-white transition-colors duration-200 text-sm font-medium px-2 py-2 rounded-lg hover:bg-white/10"
                            >
                                รายงาน
                            </a>
                        </nav>
                    </div>
                )}
            </div>
        </header>
    );
};

export default Header;
