// api.js is for axios api calls

import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

export const getItems = () => axios.get(`${API_URL}/items`);
export const getItem = (id) => axios.get(`${API_URL}/items/${id}`);
export const createItem = (item) => axios.post(`${API_URL}/items`, item);
export const updateItem = (id, item) => axios.put(`${API_URL}/items/${id}`, item);
export const deleteItem = (id) => axios.delete(`${API_URL}/items/${id}`);

export const getLists = () => axios.get(`${API_URL}/lists`);
export const getList = (id) => axios.get(`${API_URL}/lists/${id}`);
export const createList = (title) => axios.post(`${API_URL}/lists`, title);
export const updateList = (id, title) => axios.put(`${API_URL}/lists/${id}`, title);
export const deleteList = (id) => axios.delete(`${API_URL}/lists/${id}`);