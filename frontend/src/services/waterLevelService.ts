import axios from 'axios'
import type { MapMarkerResponse } from '../types/waterLevel';
import type { WaterDetailResponse } from '../types/waterDetail';

const API_BASE_URL = "http://localhost:8080";

export const getMapMarkersService = async (): Promise<MapMarkerResponse | null> => { 
    const response = await axios.get(`${API_BASE_URL}/markers`) 
    return response.data;
};

export const getMapMarkerDetailService = async(stationId: string): Promise<WaterDetailResponse | null> => {
    const response = await axios.get(`${API_BASE_URL}/markers/detail?station_id=${stationId}`)
    return response.data;
}