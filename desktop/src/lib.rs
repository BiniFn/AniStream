#[cfg(target_os = "macos")]
use tauri::TitleBarStyle;
use tauri::{Url, WebviewUrl, WebviewWindowBuilder};

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .setup(|app| {
            let url = get_webview_url();
            let win_builder = WebviewWindowBuilder::new(
                app,
                "main",
                WebviewUrl::External(Url::parse(url).unwrap()),
            )
            .title("")
            .resizable(true)
            .fullscreen(false)
            .decorations(true)
            .visible(true);

            #[cfg(target_os = "macos")]
            let win_builder = win_builder.title_bar_style(TitleBarStyle::Transparent);

            let window = win_builder.build().unwrap();

            // Set window size and center it
            if let Some(monitor) = window.current_monitor()? {
                let screen_size = monitor.size();
                let min_width = (screen_size.width as f64 * 0.75) as u32;
                let min_height = (screen_size.height as f64 * 0.75) as u32;

                window.set_min_size(Some(tauri::PhysicalSize {
                    width: min_width,
                    height: min_height,
                }))?;
                window.set_size(tauri::PhysicalSize {
                    width: min_width,
                    height: min_height,
                })?;

                window.center()?;
            }

            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

fn get_webview_url() -> &'static str {
    if cfg!(debug_assertions) {
        "http://localhost:3000"
    } else {
        "https://aniways.xyz"
    }
}
