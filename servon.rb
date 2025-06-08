class Servon < Formula
  desc "A powerful server management and development tool"
  homepage "https://github.com/angel/Servon"
  url "https://github.com/angel/Servon/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "" # This will be filled when creating a release
  license "MIT"
  head "https://github.com/angel/Servon.git", branch: "main"

  depends_on "go" => :build
  depends_on "node" => :build
  depends_on "pnpm" => :build

  def install
    # Set Go environment
    ENV["GOPATH"] = buildpath
    ENV["GO111MODULE"] = "on"
    
    # Install frontend dependencies and build
    system "pnpm", "install"
    system "pnpm", "build"
    
    # Build the Go binary
    system "go", "build", "-ldflags", "-s -w", "-o", "servon"
    
    # Install the binary
    bin.install "servon"
    
    # Install shell completions if available
    generate_completions_from_executable(bin/"servon", "completion")
  end

  test do
    # Test that the binary runs and shows version
    assert_match "servon", shell_output("#{bin}/servon --help")
    
    # Test version command if available
    system "#{bin}/servon", "version"
  end
end