import si from 'systeminformation';

export default defineEventHandler(async () => {
    const [cpu, mem, disk, os] = await Promise.all([
        si.currentLoad(),
        si.mem(),
        si.fsSize(),
        si.osInfo()
    ]);

    return {
        cpu: {
            load: cpu.currentLoad,
            cores: cpu.cpus.length
        },
        memory: {
            total: mem.total,
            used: mem.used,
            free: mem.free
        },
        disk: disk.map(d => ({
            fs: d.fs,
            size: d.size,
            used: d.used,
            mount: d.mount
        })),
        os: {
            platform: os.platform,
            distro: os.distro,
            release: os.release,
            hostname: os.hostname
        }
    };
}); 