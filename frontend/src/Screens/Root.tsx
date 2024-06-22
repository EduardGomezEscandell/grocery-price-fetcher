import React from 'react'
import Backend from "../Backend/Backend";
import LandingPage from "./LandingPage/LandingPage";
import RenderMenu from "./Menu/Menu";
import RenderPantry from "./Pantry/Pantry";
import ShoppingList from './ShoppingList/ShoppingList';
import { BrowserRouter, Route, Routes } from "react-router-dom";
import NotFound from './404/NotFound.tsx/NotFound';
import Recipes from './Recipes/Recipes';
import Products from './Products/Products';

export default function Root(): JSX.Element {
    const backend = new Backend();
    const sessionName = 'default';

    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<LandingPage backend={backend} sessionName={sessionName} />} />
                <Route path="/products" element={<Products backend={backend} sessionName={sessionName} />} />
                <Route path="/recipes" element={<Recipes backend={backend} sessionName={sessionName} />} />
                <Route path="/menu" element={<RenderMenu backend={backend} sessionName={sessionName} />} />
                <Route path="/pantry" element={<RenderPantry backend={backend} sessionName={sessionName} />} />
                <Route path="/shopping-list" element={<ShoppingList backend={backend} sessionName={sessionName} />} />
                <Route path="*" element={<NotFound/>} />
            </Routes>
        </BrowserRouter>
    );
}