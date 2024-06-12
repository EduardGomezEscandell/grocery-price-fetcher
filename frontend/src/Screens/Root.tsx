import React from 'react'
import Backend from "../Backend/Backend.ts";
import LandingPage from "./LandingPage/LandingPage.tsx";
import RenderMenu from "./Menu/Menu.tsx";
import RenderPantry from "./Pantry/Pantry.tsx";
import ShoppingList from './ShoppingList/ShoppingList.tsx';
import { BrowserRouter, Route, Routes } from "react-router-dom";

export default function Root(): JSX.Element {
    const backend = new Backend();
    const sessionName = 'default';

    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<LandingPage backend={backend} sessionName={sessionName} />} />
                <Route path="/menu" element={<RenderMenu backend={backend} sessionName={sessionName} />} />
                <Route path="/pantry" element={<RenderPantry backend={backend} sessionName={sessionName} />} />
                <Route path="/shopping-list" element={<ShoppingList backend={backend} sessionName={sessionName} />} />
            </Routes>
        </BrowserRouter>
    );
}