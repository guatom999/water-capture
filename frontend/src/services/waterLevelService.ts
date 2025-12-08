import axios from 'axios'
import type { MapMarkerResponse } from '../types/waterLevel';

const API_BASE_URL = "http://localhost:8080";

export const getMapMarkersService = async (): Promise<MapMarkerResponse> => { 
    const response = await axios.get(`${API_BASE_URL}/markers`) 
    return response.data;
};
