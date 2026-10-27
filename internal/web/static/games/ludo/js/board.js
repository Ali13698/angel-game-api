export function drawBoard() {
    const svg = document.getElementById('board');
    const NS = 'http://www.w3.org/2000/svg';
    
    // پس‌زمینه سفید
    const bg = document.createElementNS(NS, 'rect');
    bg.setAttribute('width', 550); bg.setAttribute('height', 550); bg.setAttribute('fill', '#fbffff');
    svg.appendChild(bg);

    // رسم شبکه ۱۱x۱۱
    for (let r = 0; r < 11; r++) {
        for (let c = 0; c < 11; c++) {
            const cell = document.createElementNS(NS, 'rect');
            cell.setAttribute('x', c * 50); cell.setAttribute('y', r * 50);
            cell.setAttribute('width', 50); cell.setAttribute('height', 50);
            cell.setAttribute('fill', 'none'); cell.setAttribute('stroke', '#dde5e5'); cell.setAttribute('stroke-width', '1');
            svg.appendChild(cell);
        }
    }

    // رسم پایه‌های رنگی (قرمز، سبز، آبی، زرد)
    const bases = [
        { x: 0, y: 0, c: '#e52521' },   // قرمز (بالا چپ)
        { x: 7, y: 0, c: '#1eae3a' },   // سبز (بالا راست)
        { x: 7, y: 7, c: '#007cee' },   // آبی (پایین راست)
        { x: 0, y: 7, c: '#fdc010' }    // زرد (پایین چپ)
    ];
    
    bases.forEach(b => {
        const baseRect = document.createElementNS(NS, 'rect');
        baseRect.setAttribute('x', b.x * 50); baseRect.setAttribute('y', b.y * 50);
        baseRect.setAttribute('width', 200); baseRect.setAttribute('height', 200);
        baseRect.setAttribute('fill', b.c);
        svg.appendChild(baseRect);
        
        // مربع سفید داخلی
        const inner = document.createElementNS(NS, 'rect');
        inner.setAttribute('x', b.x * 50 + 30); inner.setAttribute('y', b.y * 50 + 30);
        inner.setAttribute('width', 140); inner.setAttribute('height', 140);
        inner.setAttribute('fill', '#fff');
        svg.appendChild(inner);
    });

    // مثلث‌های مرکز
    const center = [
        { p: '200,200 200,350 275,275', c: '#e52521' },
        { p: '200,200 350,200 275,275', c: '#1eae3a' },
        { p: '350,200 350,350 275,275', c: '#007cee' },
        { p: '200,350 350,350 275,275', c: '#fdc010' }
    ];
    center.forEach(t => {
        const poly = document.createElementNS(NS, 'polygon');
        poly.setAttribute('points', t.p); poly.setAttribute('fill', t.c);
        svg.appendChild(poly);
    });
}
