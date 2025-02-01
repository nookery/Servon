import si from 'systeminformation';

export default defineEventHandler(async () => {
    try {
        const apps = await si.processes();
        return apps.list.map(app => ({
            name: app.name,
            version: app.path,
            description: app.command
        }));
    } catch (error) {
        console.error('Error fetching software information:', error);
        return [];
    }
}); 